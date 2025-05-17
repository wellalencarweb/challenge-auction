package auction

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/wellalencarweb/challenge-auction/internal/entity"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository struct {
	Db     *gorm.DB
	Logger *zap.SugaredLogger
}

func (r *Repository) CreateAuction(ctx context.Context, auction *entity.Auction) error {
	// Criação do leilão no banco de dados
	if err := r.Db.WithContext(ctx).Create(auction).Error; err != nil {
		r.Logger.Errorw("failed to create auction", "error", err)
		return err
	}

	// Obtém duração do leilão (minutos)
	duration, err := r.getAuctionDuration()
	if err != nil {
		r.Logger.Warnw("using default auction duration", "default", "5m", "error", err)
		duration = 5 * time.Minute
	}

	// Inicia goroutine para fechamento automático
	go r.scheduleAuctionClosing(ctx, auction.ID, duration)

	return nil
}

func (r *Repository) getAuctionDuration() (time.Duration, error) {
	durationStr := os.Getenv("AUCTION_DURATION_MINUTES")
	if durationStr == "" {
		return 0, strconv.ErrSyntax
	}

	minutes, err := strconv.Atoi(durationStr)
	if err != nil {
		return 0, err
	}

	return time.Duration(minutes) * time.Minute, nil
}

func (r *Repository) scheduleAuctionClosing(ctx context.Context, auctionID string, duration time.Duration) {
	select {
	case <-time.After(duration):
		r.closeAuction(ctx, auctionID)
	case <-ctx.Done():
		r.Logger.Infow("auction closing canceled", "auctionID", auctionID, "reason", ctx.Err())
	}
}

func (r *Repository) closeAuction(ctx context.Context, auctionID string) {
	result := r.Db.WithContext(ctx).
		Model(&entity.Auction{}).
		Where("id = ? AND status = 'open'", auctionID).
		Update("status", "closed")

	if result.Error != nil {
		r.Logger.Errorw("failed to close auction",
			"auctionID", auctionID,
			"error", result.Error)
	}
}
