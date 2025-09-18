package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AutoCloseManager struct {
	collection *mongo.Collection
	interval   time.Duration
	batchSize  int
	stop       chan bool
}

func NewAutoCloseManager(database *mongo.Database) *AutoCloseManager {
	interval := getCheckInterval()
	batchSize := getBatchSize()

	return &AutoCloseManager{
		collection: database.Collection("auctions"),
		interval:   interval,
		batchSize:  batchSize,
		stop:       make(chan bool),
	}
}

func (acm *AutoCloseManager) Start() {
	go func() {
		ticker := time.NewTicker(acm.interval)
		defer ticker.Stop()

		for {
			select {
			case <-acm.stop:
				return
			case <-ticker.C:
				acm.closeExpiredAuctions()
			}
		}
	}()
}

func (acm *AutoCloseManager) Stop() {
	acm.stop <- true
}

func (acm *AutoCloseManager) closeExpiredAuctions() {
	ctx := context.Background()
	now := time.Now().Unix()

	// Primeiro vamos verificar se existem leilões expirados
	// Apenas um filtro para os leilões ativos e expirados
	filter := bson.M{
		"status":   int32(0), // Active = 0
		"end_time": bson.M{"$lte": now},
	}

	logger.Info(fmt.Sprintf("Procurando leilões expirados com filtro: %v", filter))

	// Atualizando o status dos leilões expirados
	update := bson.M{
		"$set": bson.M{
			"status": int32(1), // Completed = 1
		},
	}

	logger.Info(fmt.Sprintf("Aplicando update: %v", update))

	logger.Info(fmt.Sprintf("Trying to close auctions with filter: %v", filter))
	logger.Info(fmt.Sprintf("Update to be applied: %v", update))

	result, err := acm.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error closing expired auctions", err)
		return
	}

	logger.Info(fmt.Sprintf("Update result: %+v", result))

	if result.ModifiedCount > 0 {
		logger.Info(fmt.Sprintf("Closed %d expired auctions", result.ModifiedCount))
	}
}

func getCheckInterval() time.Duration {
	defaultInterval := 30 // 30 seconds default
	if envInterval := os.Getenv("AUCTION_CHECK_INTERVAL_SECONDS"); envInterval != "" {
		if parsed, err := strconv.Atoi(envInterval); err == nil {
			defaultInterval = parsed
		}
	}
	return time.Duration(defaultInterval) * time.Second
}

func getBatchSize() int {
	defaultSize := 10
	if envSize := os.Getenv("AUCTION_BATCH_SIZE"); envSize != "" {
		if parsed, err := strconv.Atoi(envSize); err == nil {
			defaultSize = parsed
		}
	}
	return defaultSize
}
