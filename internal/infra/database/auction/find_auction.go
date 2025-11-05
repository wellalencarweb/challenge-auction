package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (ar *AuctionRepository) FindAuctionById(
	ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("Error trying to find auction by id = %s", id), err)
		return nil, internal_error.NewInternalServerError("Error trying to find auction by id")
	}

	return &auction_entity.Auction{
		Id:          auctionEntityMongo.Id,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
	}, nil
}

func (repo *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category string,
	productName string) ([]auction_entity.Auction, *internal_error.InternalError) {

	logger.Info("Iniciando busca de leilões",
		zap.Int("status_desejado", int(status)))

	// Usando bson.D para garantir a ordem dos campos e tipos corretos
	statusInt := int(status)
	logger.Info("Buscando leilões com status", zap.Int("status", statusInt))

	// Construindo o filtro usando bson.M para maior flexibilidade
	filter := bson.M{
		"status": int(status),
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["product_name"] = primitive.Regex{Pattern: productName, Options: "i"}
	}

	logger.Info("Filtro da consulta",
		zap.Any("filter", filter),
		zap.Int("status_esperado", int(status)))

	// Debug: Vamos verificar todos os documentos primeiro
	allDocs, err := repo.Collection.Find(ctx, bson.M{})
	if err != nil {
		logger.Error("Error finding all auctions", err)
		return nil, internal_error.NewInternalServerError("Error finding auctions")
	}
	defer allDocs.Close(ctx)

	var allAuctions []bson.M
	if err := allDocs.All(ctx, &allAuctions); err != nil {
		logger.Error("Error decoding all auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding auctions")
	}
	logger.Info("All documents in collection", zap.Any("auctions", allAuctions))

	// Agora vamos buscar com o filtro
	logger.Info("Executing MongoDB query with filter", zap.Any("filter", filter))
	cursor, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		logger.Error("Error finding auctions with filter", err, zap.Any("filter", filter))
		return nil, internal_error.NewInternalServerError("Error finding auctions")
	}
	defer cursor.Close(ctx)

	var auctionsMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionsMongo); err != nil {
		logger.Error("Error decoding auctions", err)
		return nil, internal_error.NewInternalServerError("Error decoding auctions")
	}

	var auctionsEntity []auction_entity.Auction
	for _, auction := range auctionsMongo {
		logger.Info(fmt.Sprintf("Found auction: ID=%s, Status=%d", auction.Id, auction.Status))
		auctionsEntity = append(auctionsEntity, auction_entity.Auction{
			Id:          auction.Id,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Status:      auction.Status,
			Description: auction.Description,
			Condition:   auction.Condition,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	}

	return auctionsEntity, nil
}
