package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestAuctionAutoClose(t *testing.T) {
	// Setup do banco de dados
	database, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Limpar a coleção antes do teste
	err := database.Collection("auctions").Drop(context.Background())
	assert.Nil(t, err)

	// Criar o repositório
	repo := NewAuctionRepository(database)

	// Criar um leilão
	auction, ierr := auction_entity.CreateAuction(
		"Produto Teste",
		"Categoria Teste",
		"Descrição completa do produto de teste",
		auction_entity.New,
	)
	assert.Nil(t, ierr)

	// Forçar o tempo de término para o passado
	auction.EndTime = time.Now().Add(-1 * time.Second)
	auction.Status = auction_entity.Active
	t.Logf("[DEBUG] Criado leilão com ID: %s", auction.Id)
	t.Logf("[DEBUG] EndTime: %v (Unix: %v)", auction.EndTime, auction.EndTime.Unix())

	// Salvar o leilão
	ierr = repo.CreateAuction(context.Background(), auction)
	assert.Nil(t, ierr)

	// Verificar estado inicial no MongoDB
	var initialState bson.M
	err = repo.Collection.FindOne(context.Background(), bson.M{"_id": auction.Id}).Decode(&initialState)
	assert.Nil(t, err)
	t.Logf("[DEBUG] Estado inicial: %+v", initialState)
	assert.Equal(t, int32(0), initialState["status"].(int32), "Status inicial deveria ser Active (0)")

	// Criar e iniciar o AutoCloseManager
	acm := NewAutoCloseManager(database)

	// Executar fechamento manual
	t.Log("[DEBUG] Executando fechamento...")
	acm.closeExpiredAuctions()
	t.Log("[DEBUG] Fechamento executado")

	// Verificar estado final no MongoDB
	var finalState bson.M
	err = repo.Collection.FindOne(context.Background(), bson.M{"_id": auction.Id}).Decode(&finalState)
	assert.Nil(t, err)
	t.Logf("[DEBUG] Estado final: %+v", finalState)

	// Verificar status
	assert.Equal(t, int32(1), finalState["status"].(int32), "Status final deveria ser Completed (1)")

	// Verificar que não há duplicação
	count, err := repo.Collection.CountDocuments(context.Background(), bson.M{"_id": auction.Id})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count, "Deveria haver exatamente um documento")
}
