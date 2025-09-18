package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"os"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AutoCloseManager struct {
	collection *mongo.Collection
	interval   time.Duration
	batchSize  int
	stop       chan struct{}
	wg         sync.WaitGroup
	mutex      sync.Mutex
	isRunning  bool
}

func NewAutoCloseManager(database *mongo.Database) *AutoCloseManager {
	interval := getCheckInterval()
	batchSize := getBatchSize()

	return &AutoCloseManager{
		collection: database.Collection("auctions"),
		interval:   interval,
		batchSize:  batchSize,
		stop:       make(chan struct{}),
		isRunning:  false,
	}
}

func (acm *AutoCloseManager) Start() {
	acm.mutex.Lock()
	if acm.isRunning {
		acm.mutex.Unlock()
		return
	}
	acm.isRunning = true
	acm.mutex.Unlock()

	acm.wg.Add(1)
	go func() {
		defer acm.wg.Done()
		ticker := time.NewTicker(acm.interval)
		defer ticker.Stop()

		for {
			select {
			case <-acm.stop:
				logger.Info("[DEBUG] AutoCloseManager: Parando...")
				return
			case <-ticker.C:
				acm.closeExpiredAuctions()
			}
		}
	}()
}

func (acm *AutoCloseManager) Stop() {
	acm.mutex.Lock()
	if !acm.isRunning {
		acm.mutex.Unlock()
		return
	}
	acm.isRunning = false
	acm.mutex.Unlock()

	close(acm.stop)
	acm.wg.Wait()
	logger.Info("[DEBUG] AutoCloseManager: Parado com sucesso")
}

func (acm *AutoCloseManager) closeExpiredAuctions() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().Unix()
	logger.Info(fmt.Sprintf("[DEBUG] AutoCloseManager: Horário atual: %v", now))

	// Buscar leilões expirados
	filter := bson.M{
		"status":   int32(0), // Active = 0
		"end_time": bson.M{"$lt": now},
	}

	logger.Info(fmt.Sprintf("[DEBUG] AutoCloseManager: Buscando leilões com filtro: %+v", filter))

	cursor, err := acm.collection.Find(ctx, filter)
	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] AutoCloseManager: Erro ao buscar leilões: %v", err))
		return
	}
	defer cursor.Close(ctx)

	var auctions []bson.M
	if err = cursor.All(ctx, &auctions); err != nil {
		logger.Error(fmt.Sprintf("[ERROR] AutoCloseManager: Erro ao decodificar leilões: %v", err))
		return
	}

	logger.Info(fmt.Sprintf("[DEBUG] AutoCloseManager: Encontrados %d leilões", len(auctions)))
	for _, a := range auctions {
		logger.Info(fmt.Sprintf("[DEBUG] AutoCloseManager: Leilão - ID: %v, Status: %v, EndTime: %v",
			a["_id"], a["status"], a["end_time"]))
	}

	if len(auctions) == 0 {
		logger.Info("[DEBUG] AutoCloseManager: Nenhum leilão para fechar")
		return
	}

	// Atualizar status
	update := bson.M{
		"$set": bson.M{
			"status": int32(1), // Completed = 1
		},
	}

	result, err := acm.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error(fmt.Sprintf("[ERROR] AutoCloseManager: Erro ao atualizar leilões: %v", err))
		return
	}

	logger.Info(fmt.Sprintf("[DEBUG] AutoCloseManager: Atualizados %d leilões", result.ModifiedCount))

	// Verificar atualizações
	for _, auction := range auctions {
		var updated bson.M
		err := acm.collection.FindOne(ctx, bson.M{"_id": auction["_id"]}).Decode(&updated)
		if err != nil {
			logger.Error(fmt.Sprintf("[ERROR] AutoCloseManager: Erro ao verificar leilão %v: %v", auction["_id"], err))
			continue
		}
		logger.Info(fmt.Sprintf("[DEBUG] AutoCloseManager: Leilão %v atualizado - Status: %v",
			updated["_id"], updated["status"]))
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

func getAuctionIds(auctions []bson.M) []interface{} {
	ids := make([]interface{}, len(auctions))
	for i, auction := range auctions {
		ids[i] = auction["_id"]
	}
	return ids
}
