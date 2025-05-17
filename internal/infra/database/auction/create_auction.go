package auction

import (
	"context"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	return nil
}

func (repo *AuctionRepository) StartAuctionMonitor(ctx context.Context) {
	durationStr := os.Getenv("AUCTION_DURATION_SECONDS")
	durationSeconds, err := strconv.Atoi(durationStr)
	if err != nil {
		logger.Error("Invalid AUCTION_DURATION_SECONDS value", err)
		return
	}

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				now := time.Now().Unix()
				filter := bson.M{
					"status":    auction_entity.OPEN,
					"timestamp": bson.M{"$lte": now - int64(durationSeconds)},
				}
				update := bson.M{"$set": bson.M{"status": auction_entity.CLOSED}}

				_, err := repo.Collection.UpdateMany(ctx, filter, update, options.Update())
				if err != nil {
					logger.Error("Error updating expired auctions", err)
				} else {
					logger.Info("Checked for expired auctions")
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}
