package auction

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/wellalencarweb/challenge-auction/internal/entity/auction_entity"
	"github.com/wellalencarweb/challenge-auction/internal/infra/database"
	"github.com/wellalencarweb/challenge-auction/internal/internal_error"
)

type CreateAuctionInputDTO struct {
	ProductName string  `json:"product_name"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type CreateAuctionOutputDTO struct {
	ID          string    `json:"id"`
	ProductName string    `json:"product_name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateAuctionUseCase struct {
	AuctionDB database.AuctionInterface
}

func NewCreateAuctionUseCase(auctionDB database.AuctionInterface) *CreateAuctionUseCase {
	return &CreateAuctionUseCase{
		AuctionDB: auctionDB,
	}
}

func (uc *CreateAuctionUseCase) Execute(input CreateAuctionInputDTO) (*CreateAuctionOutputDTO, *internal_error.InternalError) {
	auction := auction_entity.NewAuction(
		input.ProductName,
		input.Category,
		input.Description,
		input.Price,
	)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := uc.AuctionDB.Create(ctx, auction)
	if err != nil {
		return nil, internal_error.NewInternalServerError(err.Error())
	}

	go uc.startAuctionClosingRoutine(auction.ID)

	output := &CreateAuctionOutputDTO{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Price:       auction.Price,
		Status:      auction.Status,
		CreatedAt:   auction.CreatedAt,
	}

	return output, nil
}

func (uc *CreateAuctionUseCase) startAuctionClosingRoutine(auctionID string) {
	durationMinutes, _ := strconv.Atoi(os.Getenv("AUCTION_DURATION_MINUTES"))
	checkIntervalSeconds, _ := strconv.Atoi(os.Getenv("AUCTION_CHECK_INTERVAL_SECONDS"))

	if durationMinutes <= 0 {
		durationMinutes = 10 // default
	}
	if checkIntervalSeconds <= 0 {
		checkIntervalSeconds = 60 // default
	}

	auctionDuration := time.Duration(durationMinutes) * time.Minute
	checkInterval := time.Duration(checkIntervalSeconds) * time.Second

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	auctionEndTime := time.Now().Add(auctionDuration)

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			if now.After(auctionEndTime) {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()

				auction, err := uc.AuctionDB.FindById(ctx, auctionID)
				if err != nil {
					fmt.Printf("Error finding auction: %v\n", err)
					return
				}

				if auction.Status == auction_entity.StatusOpen {
					auction.Status = auction_entity.StatusClosed
					err = uc.AuctionDB.Update(ctx, auction)
					if err != nil {
						fmt.Printf("Error closing auction: %v\n", err)
					} else {
						fmt.Printf("Auction %s closed automatically\n", auctionID)
					}
				}
				return
			}
		}
	}
}
