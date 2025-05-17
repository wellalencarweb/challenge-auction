
package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAutoCloseAuction(t *testing.T) {
	os.Setenv("AUCTION_DURATION_SECONDS", "5")

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("failed to connect to MongoDB: %v", err)
	}
	db := client.Database("auction_db_test")
	db.Collection("auctions").Drop(ctx)

	repo := NewAuctionRepository(db)

	auction := &auction_entity.Auction{
		ID:          "test-auction-1",
		ProductName: "Test Product",
		Category:    "Test Category",
		Description: "Testing",
		Condition:   auction_entity.NEW,
		Status:      auction_entity.OPEN,
	}

	_, err = repo.CreateAuction(ctx, auction)
	if err != nil {
		t.Fatalf("failed to create auction: %v", err)
	}

	time.Sleep(7 * time.Second)

	var result AuctionEntityMongo
	err = db.Collection("auctions").FindOne(ctx, bson.M{"_id": auction.ID}).Decode(&result)
	if err != nil {
		t.Fatalf("failed to fetch auction: %v", err)
	}

	if result.Status != auction_entity.CLOSED {
		t.Fatalf("expected auction to be CLOSED but got %s", result.Status)
	}
}
