package auction

import (
	"context"
	"os"
	"strconv"
	"time"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
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
	auctionEnt *auction_entity.Auction,
) (*auction_entity.Auction, error) {
	auctionToInsert := AuctionEntityMongo{
		Id:          auctionEnt.Id,
		ProductName: auctionEnt.ProductName,
		Category:    auctionEnt.Category,
		Description: auctionEnt.Description,
		Condition:   auctionEnt.Condition,
		Status:      auctionEnt.Status,
		Timestamp:   time.Now().Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionToInsert)
	if err != nil {
		logger.Error("error inserting auction", err)
		return nil, internal_error.NewInternalServerError("Error inserting auction")
	}

	go ar.startAuctionTimer(auctionToInsert.Id)

	return auctionEnt, nil
}

func (ar *AuctionRepository) startAuctionTimer(auctionID string) {
	durationStr := os.Getenv("AUCTION_DURATION_SECONDS")
	durationSec, err := strconv.Atoi(durationStr)
	if err != nil {
		durationSec = 60 // valor padrão
	}

	time.Sleep(time.Duration(durationSec) * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": auctionID, "status": auction_entity.Active}
	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}

	_, err = ar.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Error("failed to auto-close auction", err)
	}
}
