package usecases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shownest/inventory-service/internal/dto/request"
	"github.com/shownest/inventory-service/internal/dto/response"
	"github.com/shownest/inventory-service/internal/mapper"
	"github.com/shownest/inventory-service/internal/repository"
	"github.com/shownest/inventory-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
)

type UseCase struct {
	repo  *repository.Repository
	cache *redis.Client
}

func New(repo *repository.Repository, cache *redis.Client) *UseCase {
	return &UseCase{repo: repo, cache: cache}
}

func (uc *UseCase) CreateSeatCategory(ctx context.Context, req request.CreateSeatCategoryRequest, hallID string) (*response.SeatCategoryInfo, error) {
	c, err := uc.repo.CreateSeatCategory(ctx, hallID, req.Name, req.PriceMultiplier)
	if err != nil {
		return nil, err
	}
	info := mapper.ToSeatCategoryInfo(c)
	return &info, nil
}

func (uc *UseCase) ListSeatCategories(ctx context.Context, hallID string) ([]response.SeatCategoryInfo, error) {
	cats, err := uc.repo.ListSeatCategoriesByHall(ctx, hallID)
	if err != nil {
		return nil, err
	}
	return mapper.ToSeatCategoryInfoList(cats), nil
}

func (uc *UseCase) BulkCreateSeats(ctx context.Context, hallID string, req request.BulkCreateSeatsRequest) ([]response.SeatInfo, error) {
	type seatInput struct {
		CategoryID string
		Row        string
		Number     int
		X, Y       *float64
	}
	inputs := make([]struct {
		CategoryID string
		Row        string
		Number     int
		X, Y       *float64
	}, len(req.Seats))
	for i, s := range req.Seats {
		inputs[i].CategoryID = s.CategoryID
		inputs[i].Row = s.Row
		inputs[i].Number = s.Number
		inputs[i].X = s.XPosition
		inputs[i].Y = s.YPosition
	}

	seats, err := uc.repo.BulkCreateSeats(ctx, hallID, inputs)
	if err != nil {
		return nil, err
	}
	return mapper.ToSeatInfoList(seats), nil
}

func (uc *UseCase) ListSeats(ctx context.Context, hallID string) ([]response.SeatInfo, error) {
	seats, err := uc.repo.ListSeatsByHall(ctx, hallID)
	if err != nil {
		return nil, err
	}
	return mapper.ToSeatInfoList(seats), nil
}

func (uc *UseCase) PublishShowtimeSeats(ctx context.Context, showtimeID, hallID string) (int, error) {
	n, err := uc.repo.PublishShowtimeSeats(ctx, showtimeID, hallID)
	if err != nil {
		return 0, err
	}
	// Bust any stale availability cache for this showtime.
	uc.cache.Del(ctx, utils.SeatAvailCacheKey(showtimeID))
	return n, nil
}

func (uc *UseCase) ListShowtimeSeats(ctx context.Context, showtimeID string) ([]response.ShowtimeSeatInfo, error) {
	cacheKey := utils.SeatAvailCacheKey(showtimeID)

	// Seat availability is high-churn (changes on every lock/release), so we use a
	// short 10-second TTL. Cache is also invalidated explicitly on every write.
	if cached, err := uc.cache.Get(ctx, cacheKey).Bytes(); err == nil {
		var infos []response.ShowtimeSeatInfo
		if json.Unmarshal(cached, &infos) == nil {
			return infos, nil
		}
	}

	ssRows, seatRows, err := uc.repo.ListShowtimeSeats(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	// Build a lookup map so we can enrich each showtime_seat row with its
	// seat metadata (row, number, category) in O(1) rather than O(n²).
	seatMap := make(map[string]struct {
		Row        string
		Number     int
		CategoryID string
	}, len(seatRows))

	for _, s := range seatRows {
		seatMap[s.ID] = struct {
			Row        string
			Number     int
			CategoryID string
		}{s.Row, s.Number, s.CategoryID}
	}

	infos := make([]response.ShowtimeSeatInfo, len(ssRows))
	for i, ss := range ssRows {
		meta := seatMap[ss.SeatID]
		infos[i] = mapper.ToShowtimeSeatInfo(&ss, meta.Row, meta.Number, meta.CategoryID)
	}

	if b, err := json.Marshal(infos); err == nil {
		ttl := time.Duration(utils.SeatAvailCacheTTL) * time.Second
		if err := uc.cache.Set(ctx, cacheKey, b, ttl).Err(); err != nil {
			logger.Get().Warn("cache set seat availability failed",
				zap.String("showtimeId", showtimeID), zap.Error(err))
		}
	}
	return infos, nil
}

func (uc *UseCase) LockSeats(ctx context.Context, showtimeID, userID string, req request.LockSeatsRequest) (*response.LockSeatsResponse, error) {
	locked := make([]string, 0, len(req.SeatIDs))
	failed := make([]string, 0)

	lockTTL := time.Duration(utils.SeatLockTTLSeconds) * time.Second

	for _, seatID := range req.SeatIDs {
		key := utils.SeatLockKey(showtimeID, seatID)

		// SET NX is atomic: sets the key only if it does not already exist.
		// Returns "OK" on success, empty string if the key was already held.
		// This is the first line of defence against concurrent seat selection.
		res, err := uc.cache.SetArgs(ctx, key, userID, redis.SetArgs{TTL: lockTTL, Mode: "NX"}).Result()
		if err != nil {
			return nil, apperrors.Wrap(apperrors.CodeInternal, "redis lock seat", err)
		}
		if res != "OK" {
			// Another user already holds the Redis lock — reject immediately without hitting the DB.
			failed = append(failed, seatID)
			continue
		}

		// Redis lock acquired. Now persist the lock in the DB via a conditional UPDATE
		// (WHERE status = 'available') to guard against any edge case where Redis
		// and DB get out of sync.
		tx, err := uc.repo.BeginTx(ctx)
		if err != nil {
			// Roll back the Redis key so the seat doesn't stay locked indefinitely.
			uc.cache.Del(ctx, key)
			return nil, err
		}

		if err := uc.repo.LockShowtimeSeat(ctx, tx, showtimeID, seatID, userID); err != nil {
			// DB says seat is no longer available — release the Redis key and mark as failed.
			tx.Rollback(ctx)
			uc.cache.Del(ctx, key)
			failed = append(failed, seatID)
			continue
		}

		if err := tx.Commit(ctx); err != nil {
			uc.cache.Del(ctx, key)
			return nil, apperrors.Wrap(apperrors.CodeDBError, "commit lock transaction", err)
		}

		locked = append(locked, seatID)
	}

	// Invalidate availability cache so the next read reflects the new lock statuses.
	uc.cache.Del(ctx, utils.SeatAvailCacheKey(showtimeID))

	return &response.LockSeatsResponse{Locked: locked, Failed: failed}, nil
}

func (uc *UseCase) ReleaseSeats(ctx context.Context, showtimeID, userID string, req request.ReleaseSeatsRequest) error {
	for _, seatID := range req.SeatIDs {
		key := utils.SeatLockKey(showtimeID, seatID)

		// Only delete the Redis key if it belongs to this user.
		// This prevents a user from releasing a seat locked by someone else.
		val, err := uc.cache.Get(ctx, key).Result()
		if err == nil && val == userID {
			uc.cache.Del(ctx, key)
		}

		if err := uc.repo.ReleaseShowtimeSeat(ctx, showtimeID, seatID, userID); err != nil {
			return err
		}
	}

	uc.cache.Del(ctx, utils.SeatAvailCacheKey(showtimeID))
	return nil
}

func (uc *UseCase) ConfirmSeats(ctx context.Context, showtimeID, userID string, req request.ConfirmSeatsRequest) error {
	if err := uc.repo.BookShowtimeSeats(ctx, showtimeID, userID, req.SeatIDs); err != nil {
		return err
	}

	// Remove the Redis lock keys — seats are now permanently booked, not just locked.
	for _, seatID := range req.SeatIDs {
		uc.cache.Del(ctx, utils.SeatLockKey(showtimeID, seatID))
	}

	uc.cache.Del(ctx, utils.SeatAvailCacheKey(showtimeID))
	return nil
}
