package auction

import (
	"context"
	"fmt"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	testDatabase *mongo.Database
	testClient   *mongo.Client
)

func setupIntegrationTest(t *testing.T) (*mongo.Database, func()) {
	if testDatabase != nil {
		// Limpar a coleção antes do teste
		testDatabase.Collection("auctions").Drop(context.Background())
		return testDatabase, func() {}
	}

	// Configurar conexão com MongoDB
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://admin:admin@localhost:27017"
	}

	// Conectar ao MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping para garantir que a conexão está OK
	err = client.Ping(context.Background(), nil)
	if err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully")

	// Usar banco de dados de teste
	database := client.Database("auction_test")
	testDatabase = database
	testClient = client

	// Retornar função cleanup
	return database, func() {
		if testClient != nil {
			testClient.Disconnect(context.Background())
		}
	}
}
