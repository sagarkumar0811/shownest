package usecases

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shownest/catalog-service/internal/config"
	"github.com/shownest/catalog-service/internal/dto/request"
	"github.com/shownest/catalog-service/internal/dto/response"
	"github.com/shownest/catalog-service/internal/mapper"
	"github.com/shownest/catalog-service/internal/repository"
	"github.com/shownest/catalog-service/internal/utils"
	pkgaws "github.com/shownest/pkg/aws"
	apperrors "github.com/shownest/pkg/errors"
	"go.uber.org/zap"

	"github.com/shownest/pkg/logger"
)

const (
	eventCacheTTL    = 5 * time.Minute
	showtimeCacheTTL = 1 * time.Minute
)

type UseCase struct {
	repo   *repository.Repository
	s3     *pkgaws.S3Client
	cache  *redis.Client
	config *config.Config
}

func New(repo *repository.Repository, s3 *pkgaws.S3Client, cache *redis.Client, cfg *config.Config) *UseCase {
	return &UseCase{repo: repo, s3: s3, cache: cache, config: cfg}
}

func (uc *UseCase) CreateEvent(ctx context.Context, userID string, req request.CreateEventRequest) (*response.EventInfo, error) {
	e, err := uc.repo.CreateEvent(ctx,
		userID, req.MerchantID, req.Title, req.Description,
		req.Category, req.Language, req.DurationMinutes, req.Rating,
	)
	if err != nil {
		return nil, err
	}
	info := mapper.ToEventInfo(e)
	return &info, nil
}

func (uc *UseCase) GetEvent(ctx context.Context, eventID string) (*response.EventInfo, error) {
	cacheKey := utils.EventCacheKey(eventID)

	if cached, err := uc.cache.Get(ctx, cacheKey).Bytes(); err == nil {
		var info response.EventInfo
		if json.Unmarshal(cached, &info) == nil {
			return &info, nil
		}
	}

	e, err := uc.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	info := mapper.ToEventInfo(e)

	if b, err := json.Marshal(info); err == nil {
		if err := uc.cache.Set(ctx, cacheKey, b, eventCacheTTL).Err(); err != nil {
			logger.Get().Warn("cache set event failed", zap.String("eventId", eventID), zap.Error(err))
		}
	}
	return &info, nil
}

func (uc *UseCase) ListEvents(ctx context.Context, req request.ListEventsRequest) ([]response.EventInfo, error) {
	events, err := uc.repo.ListEvents(ctx, req.Category, req.MerchantID, req.Status)
	if err != nil {
		return nil, err
	}
	return mapper.ToEventInfoList(events), nil
}

func (uc *UseCase) UpdateEvent(ctx context.Context, userID, eventID string, req request.UpdateEventRequest) (*response.EventInfo, error) {
	existing, err := uc.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if existing.UserID != userID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "event does not belong to your merchant account")
	}

	fields := map[string]any{}
	if req.Title != nil {
		fields["title"] = *req.Title
	}
	if req.Description != nil {
		fields["description"] = *req.Description
	}
	if req.Language != nil {
		fields["language"] = *req.Language
	}
	if req.DurationMinutes != nil {
		fields["duration_minutes"] = *req.DurationMinutes
	}
	if req.Rating != nil {
		fields["rating"] = *req.Rating
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if len(fields) == 0 {
		info := mapper.ToEventInfo(existing)
		return &info, nil
	}

	e, err := uc.repo.UpdateEvent(ctx, eventID, fields)
	if err != nil {
		return nil, err
	}

	uc.cache.Del(ctx, utils.EventCacheKey(eventID))

	info := mapper.ToEventInfo(e)
	return &info, nil
}

func (uc *UseCase) CreateShowtime(ctx context.Context, userID, eventID string, req request.CreateShowtimeRequest) (*response.ShowtimeInfo, error) {
	event, err := uc.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event.UserID != userID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "event does not belong to your merchant account")
	}

	expectedEnd := req.StartTime.Add(time.Duration(event.DurationMinutes) * time.Minute)
	if !req.EndTime.Equal(expectedEnd) {
		return nil, apperrors.New(apperrors.CodeInvalidArgument, "endTime must equal startTime + event duration")
	}

	tx, err := uc.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	conflict, err := uc.repo.CheckHallConflict(ctx, tx, req.HallID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, apperrors.New(apperrors.CodeAlreadyExists, "hall already has a showtime in this time window")
	}

	s, err := uc.repo.CreateShowtimeInTx(ctx, tx, eventID, req.HallID, req.StartTime, req.EndTime, req.BasePrice)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "commit transaction", err)
	}

	uc.cache.Del(ctx, utils.ShowtimeListCacheKey(eventID))

	info := mapper.ToShowtimeInfo(s)
	return &info, nil
}

func (uc *UseCase) GetShowtime(ctx context.Context, showtimeID string) (*response.ShowtimeInfo, error) {
	s, err := uc.repo.GetShowtimeByID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}
	info := mapper.ToShowtimeInfo(s)
	return &info, nil
}

func (uc *UseCase) GetShowtimeBasePrice(ctx context.Context, showtimeID string) (float64, error) {
	s, err := uc.repo.GetShowtimeByID(ctx, showtimeID)
	if err != nil {
		return 0, err
	}
	price, err := s.BasePrice.Float64Value()
	if err != nil || !price.Valid {
		return 0, apperrors.New(apperrors.CodeInternal, "invalid base price on showtime")
	}
	return price.Float64, nil
}

func (uc *UseCase) ListShowtimes(ctx context.Context, eventID string) ([]response.ShowtimeInfo, error) {
	cacheKey := utils.ShowtimeListCacheKey(eventID)

	if cached, err := uc.cache.Get(ctx, cacheKey).Bytes(); err == nil {
		var infos []response.ShowtimeInfo
		if json.Unmarshal(cached, &infos) == nil {
			return infos, nil
		}
	}

	showtimes, err := uc.repo.ListShowtimesByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	infos := mapper.ToShowtimeInfoList(showtimes)

	if b, err := json.Marshal(infos); err == nil {
		if err := uc.cache.Set(ctx, cacheKey, b, showtimeCacheTTL).Err(); err != nil {
			logger.Get().Warn("cache set showtimes failed", zap.String("eventId", eventID), zap.Error(err))
		}
	}
	return infos, nil
}

func (uc *UseCase) UpdateShowtime(ctx context.Context, userID, showtimeID string, req request.UpdateShowtimeRequest) (*response.ShowtimeInfo, error) {
	s, err := uc.repo.GetShowtimeByID(ctx, showtimeID)
	if err != nil {
		return nil, err
	}

	event, err := uc.repo.GetEventByID(ctx, s.EventID)
	if err != nil {
		return nil, err
	}
	if event.UserID != userID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "showtime does not belong to your merchant account")
	}

	updated, err := uc.repo.UpdateShowtimeStatus(ctx, showtimeID, *req.Status)
	if err != nil {
		return nil, err
	}

	uc.cache.Del(ctx, utils.ShowtimeListCacheKey(s.EventID))

	info := mapper.ToShowtimeInfo(updated)
	return &info, nil
}

func (uc *UseCase) RequestMediaUploadURL(ctx context.Context, userID, eventID, mediaType string) (*response.MediaUploadURLResponse, error) {
	event, err := uc.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event.UserID != userID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "event does not belong to your merchant account")
	}

	id, err := utils.NewUUID()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate key id", err)
	}

	s3Key := utils.GetS3Key(uc.config.App, event.MerchantID, "media", eventID, mediaType, id)
	ttl := time.Duration(utils.MediaUploadURLTTL) * time.Minute
	uploadURL, err := uc.s3.PresignPutURL(ctx, s3Key, ttl)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "generate upload url", err)
	}
	return &response.MediaUploadURLResponse{
		UploadURL: uploadURL,
		S3Key:     s3Key,
		MediaType: mediaType,
	}, nil
}

func (uc *UseCase) ConfirmMedia(ctx context.Context, userID, eventID string, req request.ConfirmMediaRequest) (*response.MediaInfo, error) {
	event, err := uc.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event.UserID != userID {
		return nil, apperrors.New(apperrors.CodePermissionDenied, "event does not belong to your merchant account")
	}

	m, err := uc.repo.CreateMedia(ctx, eventID, req.MediaType, req.S3Key)
	if err != nil {
		return nil, err
	}
	info := mapper.ToMediaInfo(m)
	return &info, nil
}

func (uc *UseCase) ListMedia(ctx context.Context, eventID string) ([]response.MediaInfo, error) {
	if _, err := uc.repo.GetEventByID(ctx, eventID); err != nil {
		return nil, err
	}
	media, err := uc.repo.GetMediaByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return mapper.ToMediaInfoList(media), nil
}
