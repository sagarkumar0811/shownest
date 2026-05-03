package jobs

import (
	"context"
	"time"

	"github.com/shownest/booking-service/internal/client"
	"github.com/shownest/booking-service/internal/repository"
	"github.com/shownest/booking-service/internal/utils"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
)

// StartExpireBookingsJob starts a background job that periodically cancels pending bookings that have expired
// and releases their reserved seats back to inventory.
func StartExpireBookingsJob(ctx context.Context, repo *repository.Repository, inventory *client.InventoryClient) {
	interval := time.Duration(utils.BookingExpirySeconds) * time.Second / 2

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				bookings, err := repo.FetchAndCancelExpiredBookings(ctx)
				if err != nil {
					logger.Get().Error("expire bookings failed", zap.Error(err))
					continue
				}

				for _, b := range bookings {
					items, err := repo.GetBookingItems(ctx, b.ID)
					if err != nil {
						logger.Get().Error("get items for expired booking failed",
							zap.String("bookingId", b.ID), zap.Error(err))
						continue
					}

					seatIDs := make([]string, len(items))
					for i, item := range items {
						seatIDs[i] = item.SeatID
					}

					if err := inventory.ReleaseSeats(ctx, b.ShowtimeID, seatIDs); err != nil {
						logger.Get().Error("release seats for expired booking failed",
							zap.String("bookingId", b.ID), zap.Error(err))
					}
				}

				if len(bookings) > 0 {
					logger.Get().Info("expired pending bookings cancelled", zap.Int("count", len(bookings)))
				}
			}
		}
	}()
}
