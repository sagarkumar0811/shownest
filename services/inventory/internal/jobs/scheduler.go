package jobs

import (
	"context"
	"time"

	"github.com/shownest/inventory-service/internal/repository"
	"github.com/shownest/inventory-service/internal/utils"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
)

func StartJobs(ctx context.Context, repo *repository.Repository) {
	StartExpireLocksJob(ctx, repo)
}

// StartExpireLocksJob starts a background job that periodically marks locked seats as available
// if they have been locked for longer than the TTL.
func StartExpireLocksJob(ctx context.Context, repo *repository.Repository) {
	interval := time.Duration(utils.SeatLockTTLSeconds) * time.Second / 2

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				n, err := repo.ExpireLockedSeats(ctx, utils.SeatLockTTLSeconds)
				if err != nil {
					logger.Get().Error("expire locked seats failed", zap.Error(err))
				} else if n > 0 {
					logger.Get().Info("expired stale seat locks", zap.Int64("count", n))
				}
			}
		}
	}()
}
