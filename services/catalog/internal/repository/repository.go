package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shownest/catalog-service/internal/models"
	"github.com/shownest/catalog-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var eventColumns = []string{
	"id", "user_id", "merchant_id", "title", "description", "category",
	"language", "duration_minutes", "rating", "status", "created_at", "updated_at", "deleted_at",
}

var showtimeColumns = []string{
	"id", "event_id", "hall_id", "start_time", "end_time", "base_price",
	"status", "created_at", "updated_at", "deleted_at",
}

var mediaColumns = []string{
	"id", "event_id", "media_type", "s3_key", "cdn_url", "created_at",
}

func (r *Repository) CreateEvent(ctx context.Context, userID, merchantID, title, description, category, language string, durationMins int, rating string) (*models.Event, error) {
	ins := psql.Insert("events").
		Columns("user_id", "merchant_id", "title", "description", "category", "language", "duration_minutes", "rating").
		Values(
			userID, merchantID, title,
			nullStr(description), category, nullStr(language),
			durationMins, nullStr(rating),
		).
		Suffix("RETURNING " + utils.JoinColumns(eventColumns))

	sql, args, err := ins.ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Event])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create event", err)
	}
	return &e, nil
}

func (r *Repository) GetEventByID(ctx context.Context, id string) (*models.Event, error) {
	sql, args, err := psql.Select(eventColumns...).From("events").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var e models.Event
	if err := pgxscan.Get(ctx, r.db, &e, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "event not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get event by id", err)
	}
	return &e, nil
}

func (r *Repository) ListEvents(ctx context.Context, category, merchantID, status string) ([]models.Event, error) {
	q := psql.Select(eventColumns...).From("events").Where("deleted_at IS NULL").OrderBy("created_at DESC")
	if category != "" {
		q = q.Where(sq.Eq{"category": category})
	}
	if merchantID != "" {
		q = q.Where(sq.Eq{"merchant_id": merchantID})
	}
	if status != "" {
		q = q.Where(sq.Eq{"status": status})
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var events []models.Event
	if err := pgxscan.Select(ctx, r.db, &events, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list events", err)
	}
	return events, nil
}

func (r *Repository) UpdateEvent(ctx context.Context, id string, fields map[string]interface{}) (*models.Event, error) {
	q := psql.Update("events").Where(sq.Eq{"id": id}).Where("deleted_at IS NULL")
	for k, v := range fields {
		q = q.Set(k, v)
	}
	q = q.Suffix("RETURNING " + utils.JoinColumns(eventColumns))

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	e, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Event])
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "event not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "update event", err)
	}
	return &e, nil
}

func (r *Repository) CheckHallConflict(ctx context.Context, tx pgx.Tx, hallID string, startTime, endTime time.Time) (bool, error) {
	sql, args, err := psql.Select("id").From("showtimes").
		Where(sq.Eq{"hall_id": hallID}).
		Where("deleted_at IS NULL").
		Where(sq.NotEq{"status": utils.ShowtimeStatusCancelled}).
		Where("start_time < ?", endTime).
		Where("end_time > ?", startTime).
		Suffix("FOR UPDATE").
		ToSql()
	if err != nil {
		return false, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var ids []string
	if err := pgxscan.Select(ctx, tx, &ids, sql, args...); err != nil {
		return false, apperrors.Wrap(apperrors.CodeDBError, "check hall conflict", err)
	}
	return len(ids) > 0, nil
}

func (r *Repository) CreateShowtime(ctx context.Context, eventID, hallID string, startTime, endTime time.Time, basePrice float64) (*models.Showtime, error) {
	return r.createShowtimeInTx(ctx, nil, eventID, hallID, startTime, endTime, basePrice)
}

func (r *Repository) CreateShowtimeInTx(ctx context.Context, tx pgx.Tx, eventID, hallID string, startTime, endTime time.Time, basePrice float64) (*models.Showtime, error) {
	return r.createShowtimeInTx(ctx, tx, eventID, hallID, startTime, endTime, basePrice)
}

func (r *Repository) createShowtimeInTx(ctx context.Context, q pgx.Tx, eventID, hallID string, startTime, endTime time.Time, basePrice float64) (*models.Showtime, error) {
	sqlStr, args, err := psql.Insert("showtimes").
		Columns("event_id", "hall_id", "start_time", "end_time", "base_price").
		Values(eventID, hallID, startTime, endTime, basePrice).
		Suffix("RETURNING " + utils.JoinColumns(showtimeColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var rows pgx.Rows
	if q != nil {
		rows, _ = q.Query(ctx, sqlStr, args...)
	} else {
		rows, _ = r.db.Query(ctx, sqlStr, args...)
	}
	s, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Showtime])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create showtime", err)
	}
	return &s, nil
}

func (r *Repository) GetShowtimeByID(ctx context.Context, id string) (*models.Showtime, error) {
	sql, args, err := psql.Select(showtimeColumns...).From("showtimes").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var s models.Showtime
	if err := pgxscan.Get(ctx, r.db, &s, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "showtime not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get showtime by id", err)
	}
	return &s, nil
}

func (r *Repository) ListShowtimesByEventID(ctx context.Context, eventID string) ([]models.Showtime, error) {
	sql, args, err := psql.Select(showtimeColumns...).From("showtimes").
		Where(sq.Eq{"event_id": eventID}).
		Where("deleted_at IS NULL").
		OrderBy("start_time ASC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var showtimes []models.Showtime
	if err := pgxscan.Select(ctx, r.db, &showtimes, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list showtimes", err)
	}
	return showtimes, nil
}

func (r *Repository) UpdateShowtimeStatus(ctx context.Context, id, status string) (*models.Showtime, error) {
	sql, args, err := psql.Update("showtimes").
		Set("status", status).
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		Suffix("RETURNING " + utils.JoinColumns(showtimeColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	s, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Showtime])
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "showtime not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "update showtime status", err)
	}
	return &s, nil
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "begin transaction", err)
	}
	return tx, nil
}

func (r *Repository) CreateMedia(ctx context.Context, eventID, mediaType, s3Key string) (*models.EventMedia, error) {
	sql, args, err := psql.Insert("event_media").
		Columns("event_id", "media_type", "s3_key").
		Values(eventID, mediaType, s3Key).
		Suffix("RETURNING " + utils.JoinColumns(mediaColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.EventMedia])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create media", err)
	}
	return &m, nil
}

func (r *Repository) GetMediaByEventID(ctx context.Context, eventID string) ([]models.EventMedia, error) {
	sql, args, err := psql.Select(mediaColumns...).From("event_media").
		Where(sq.Eq{"event_id": eventID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var media []models.EventMedia
	if err := pgxscan.Select(ctx, r.db, &media, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get media by event id", err)
	}
	return media, nil
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}
