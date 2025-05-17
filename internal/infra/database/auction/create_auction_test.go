
package auction_test

import (
    "context"
    "os"
    "testing"
    "time"

    "fullcycle-auction_go/internal/entity/auction_entity"
    auctiondb "fullcycle-auction_go/internal/infra/database/auction"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
)

func TestAuctionAutoClose(t *testing.T) {
    os.Setenv("AUCTION_DURATION_SECONDS", "2")

    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        t.Fatalf("erro ao conectar no MongoDB: %v", err)
    }
    defer client.Disconnect(ctx)

    db := client.Database("auction")
    repo := auctiondb.NewAuctionRepository(db)

    auction := &auction_entity.AuctionEntity{
        Id:          "test-auto-close",
        ProductName: "Produto Teste",
        Category:    "Categoria X",
        Description: "Teste automático",
        Condition:   auction_entity.New,
        Status:      auction_entity.Opened,
        Timestamp:   time.Now().Unix(),
    }

    err = repo.CreateAuction(ctx, auction)
    if err != nil {
        t.Fatalf("erro ao criar auction: %v", err)
    }

    time.Sleep(3 * time.Second)

    var result bson.M
    err = db.Collection("auctions").FindOne(ctx, bson.M{"_id": "test-auto-close"}).Decode(&result)
    if err != nil {
        t.Fatalf("erro ao buscar auction: %v", err)
    }

    if result["status"] != "CLOSED" {
        t.Errorf("esperava status CLOSED, obteve %v", result["status"])
    }
}
