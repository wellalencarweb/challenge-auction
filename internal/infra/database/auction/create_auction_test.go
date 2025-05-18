
package auction_test

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/infra/database/auction"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func setupTestDB() *mongo.Database {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return client.Database("auction_test_db")
}

func TestAutomaticAuctionClosure(t *testing.T) {
	db := setupTestDB()
	repo := auction.NewAuctionRepository(db)
	os.Setenv("AUCTION_DURATION_SECONDS", "2")

	a := &auction_entity.AuctionEntity{
		Id:          "test-auction-123",
		ProductName: "Produto Teste",
		Category:    "Teste",
		Description: "Descrição",
		Condition:   auction_entity.NEW,
		Status:      auction_entity.OPEN,
		Timestamp:   time.Now().Unix(),
	}

	ctx := context.Background()
	_, err := repo.CreateAuction(ctx, a)
	if err != nil {
		t.Fatalf("erro ao criar leilão: %v", err)
	}

	time.Sleep(3 * time.Second) // Aguarda o fechamento automático

	var result map[string]interface{}
	err = repo.Collection.FindOne(ctx, bson.M{"_id": a.Id}).Decode(&result)
	if err != nil {
		t.Fatalf("erro ao buscar leilão: %v", err)
	}

	if result["status"] != "closed" {
		t.Errorf("esperado status 'closed', recebido '%v'", result["status"])
	}
}
