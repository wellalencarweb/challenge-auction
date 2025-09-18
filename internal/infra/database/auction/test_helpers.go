package auction

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

// TestMongoDB encapsula o ambiente de teste do MongoDB
type TestMongoDB struct {
	Mt       *mtest.T
	Database *mongo.Database
}

// setupMongoTest configura um ambiente de teste MongoDB usando mtest
func setupMongoTest(t *testing.T, testFunc func(*mtest.T)) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test", func(mt *mtest.T) {
		testFunc(mt)
	})
}
