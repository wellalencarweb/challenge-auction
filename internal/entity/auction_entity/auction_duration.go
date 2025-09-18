package auction_entity

import (
	"os"
	"strconv"
	"time"
)

func getAuctionDuration() time.Duration {
	durationSeconds := 300 // default 5 minutes
	if envDuration := os.Getenv("AUCTION_DURATION_SECONDS"); envDuration != "" {
		if parsed, err := strconv.Atoi(envDuration); err == nil {
			durationSeconds = parsed
		}
	}
	return time.Duration(durationSeconds) * time.Second
}
