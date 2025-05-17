package auction

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wellalencarweb/challenge-auction/internal/entity/auction_entity"
	"github.com/wellalencarweb/challenge-auction/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuctionDBMock struct {
	auctions map[string]*auction_entity.Auction
}

func NewAuctionDBMock() *AuctionDBMock {
	return &AuctionDBMock{
		auctions: make(map[string]*auction_entity.Auction),
	}
}

func (m *AuctionDBMock) Create(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	if auction.ID == "" {
		auction.ID = primitive.NewObjectID().Hex()
	}
	m.auctions[auction.ID] = auction
	return nil
}

func (m *AuctionDBMock) FindById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	auction, exists := m.auctions[id]
	if !exists {
		return nil, internal_error.NewNotFoundError("auction not found")
	}
	return auction, nil
}

func (m *AuctionDBMock) Update(ctx context.Context, auction *auction_entity.Auction) *internal_error.InternalError {
	if _, exists := m.auctions[auction.ID]; !exists {
		return internal_error.NewNotFoundError("auction not found")
	}
	m.auctions[auction.ID] = auction
	return nil
}

func TestCreateAuctionAndAutoClose(t *testing.T) {
	t.Setenv("AUCTION_DURATION_MINUTES", "1")       // 1 minuto para teste
	t.Setenv("AUCTION_CHECK_INTERVAL_SECONDS", "1") // verificar a cada 1 segundo

	mockDB := NewAuctionDBMock()
	uc := NewCreateAuctionUseCase(mockDB)

	input := CreateAuctionInputDTO{
		ProductName: "Test Product",
		Category:    "Test",
		Description: "Test Description",
		Price:       100.0,
	}

	output, err := uc.Execute(input)
	assert.Nil(t, err)
	assert.Equal(t, auction_entity.StatusOpen, output.Status)

	ctx := context.Background()
	auction, err := mockDB.FindById(ctx, output.ID)
	assert.Nil(t, err)
	assert.Equal(t, auction_entity.StatusOpen, auction.Status)

	// Aguardar tempo suficiente para o leilão fechar
	time.Sleep(70 * time.Second) // mais que 1 minuto

	// Verificar se o leilão foi fechado
	auction, err = mockDB.FindById(ctx, output.ID)
	assert.Nil(t, err)
	assert.Equal(t, auction_entity.StatusClosed, auction.Status)
}
