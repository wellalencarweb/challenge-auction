package auction_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/wellalencarweb/challenge-auction/internal/infra/database/auction"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionAutoClose(t *testing.T) {
	os.Setenv("AUCTION_DURATION_SECONDS", "2")
	ctx := context.Background()

	// Conecta ao Mongo local conforme docker-compose padrão
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:root@localhost:27017"))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Disconnect(ctx)

	coll := client.Database("auction").Collection("auctions")
	repo := auction.NewAuctionRepository(coll)

	// Limpa antes
	coll.DeleteMany(ctx, bson.M{})

	// Cria um leilão com timestamp no passado
	auctionDoc := bson.M{
		"_id":          "test-auction-123",
		"product_name": "Test Product",
		"category":     "Test",
		"description":  "A test auction",
		"condition":    "NEW",
		"status":       "OPEN",
		"timestamp":    time.Now().Unix() - 10,
	}
	_, err = coll.InsertOne(ctx, auctionDoc)
	if err != nil {
		t.Fatal(err)
	}

	repo.StartAuctionMonitor(ctx)

	// Aguarda tempo suficiente para o monitor rodar
	time.Sleep(5 * time.Second)

	// Verifica se foi fechado
	var result bson.M
	err = coll.FindOne(ctx, bson.M{"_id": "test-auction-123"}).Decode(&result)
	if err != nil {
		t.Fatal(err)
	}

	if result["status"] != "CLOSED" {
		t.Errorf("Auction was not closed: got %v", result["status"])
	}
}
