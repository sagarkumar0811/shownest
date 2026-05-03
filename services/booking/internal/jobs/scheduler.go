package jobs

import (
	"context"

	"github.com/shownest/booking-service/internal/client"
	"github.com/shownest/booking-service/internal/repository"
)

func StartJobs(ctx context.Context, repo *repository.Repository, inventory *client.InventoryClient) {
	StartExpireBookingsJob(ctx, repo, inventory)
}
