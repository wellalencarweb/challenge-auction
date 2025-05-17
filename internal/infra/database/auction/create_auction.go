package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}
type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {
	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	
	// Lógica de fechamento automático usando goroutine
	go func(auctionId string) {
		durationStr := os.Getenv("AUCTION_DURATION_SECONDS")
		duration, err := strconv.Atoi(durationStr)
		if err != nil || duration <= 0 {
			duration = 30 // valor padrão de fallback (30 segundos)
		}
		time.Sleep(time.Duration(duration) * time.Second)

		filter := bson.M{"_id": auctionId, "status": "OPENED"}
		update := bson.M{"$set": bson.M{"status": "CLOSED"}}

		_, err = ar.Collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			logger.Error("erro ao fechar o leilão automaticamente", err)
		}
	}(auctionEntity.Id)

	return nil
}
