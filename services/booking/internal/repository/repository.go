package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shownest/booking-service/internal/models"
	"github.com/shownest/booking-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
	pkgutils "github.com/shownest/pkg/utils"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var bookingColumns = []string{
	"id", "user_id", "showtime_id", "status", "total_amount", "qr_token", "used_at", "expires_at", "created_at", "updated_at",
}

var bookingItemColumns = []string{
	"id", "booking_id", "seat_id", "category_id", "price", "created_at",
}

func (r *Repository) CreateBookingWithItems(ctx context.Context, userID, showtimeID string, totalAmount float64, expiresAt time.Time, items []struct {
	SeatID, CategoryID string
	Price              float64
}) (*models.Booking, []models.BookingItem, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeDBError, "begin transaction", err)
	}
	defer tx.Rollback(ctx)

	bookingSql, bookingArgs, err := psql.Insert("bookings").
		Columns("user_id", "showtime_id", "status", "total_amount", "expires_at").
		Values(userID, showtimeID, utils.BookingStatusPending, totalAmount, expiresAt).
		Suffix("RETURNING " + pkgutils.JoinColumns(bookingColumns)).
		ToSql()
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeInternal, "build booking query", err)
	}

	rows, _ := tx.Query(ctx, bookingSql, bookingArgs...)
	booking, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Booking])
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeDBError, "insert booking", err)
	}

	q := psql.Insert("booking_items").
		Columns("booking_id", "seat_id", "category_id", "price")
	for _, item := range items {
		q = q.Values(booking.ID, item.SeatID, item.CategoryID, item.Price)
	}
	q = q.Suffix("RETURNING " + pkgutils.JoinColumns(bookingItemColumns))

	itemSql, itemArgs, err := q.ToSql()
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeInternal, "build booking items query", err)
	}

	itemRows, _ := tx.Query(ctx, itemSql, itemArgs...)
	bookingItems, err := pgx.CollectRows(itemRows, pgx.RowToStructByName[models.BookingItem])
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeDBError, "insert booking items", err)
	}

	logSql, logArgs, err := psql.Insert("bookings_state_log").
		Columns("booking_id", "from_status", "to_status").
		Values(booking.ID, nil, utils.BookingStatusPending).
		ToSql()
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeInternal, "build state log query", err)
	}
	if _, err := tx.Exec(ctx, logSql, logArgs...); err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeDBError, "insert state log", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeDBError, "commit create booking", err)
	}

	return &booking, bookingItems, nil
}

func (r *Repository) GetBookingByID(ctx context.Context, id string) (*models.Booking, error) {
	sql, args, err := psql.Select(bookingColumns...).From("bookings").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var b models.Booking
	if err := pgxscan.Get(ctx, r.db, &b, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "booking not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get booking by id", err)
	}
	return &b, nil
}

func (r *Repository) GetBookingItems(ctx context.Context, bookingID string) ([]models.BookingItem, error) {
	sql, args, err := psql.Select(bookingItemColumns...).From("booking_items").
		Where(sq.Eq{"booking_id": bookingID}).
		OrderBy("created_at ASC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var items []models.BookingItem
	if err := pgxscan.Select(ctx, r.db, &items, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get booking items", err)
	}
	return items, nil
}

func (r *Repository) UpdateBookingStatus(ctx context.Context, id, status string, qrToken *string) (*models.Booking, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "begin transaction", err)
	}
	defer tx.Rollback(ctx)

	// Fetch current status for the audit log.
	var current models.Booking
	fetchSql, fetchArgs, _ := psql.Select("status").From("bookings").Where(sq.Eq{"id": id}).ToSql()
	if err := pgxscan.Get(ctx, tx, &current, fetchSql, fetchArgs...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "booking not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "fetch booking status", err)
	}

	q := psql.Update("bookings").
		Set("status", status).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING " + pkgutils.JoinColumns(bookingColumns))

	if qrToken != nil {
		q = q.Set("qr_token", *qrToken)
	}

	updateSql, updateArgs, err := q.ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build update query", err)
	}

	rows, _ := tx.Query(ctx, updateSql, updateArgs...)
	b, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Booking])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "update booking status", err)
	}

	logSql, logArgs, err := psql.Insert("bookings_state_log").
		Columns("booking_id", "from_status", "to_status").
		Values(id, current.Status, status).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build state log query", err)
	}
	if _, err := tx.Exec(ctx, logSql, logArgs...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "insert state log", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "commit update booking status", err)
	}

	return &b, nil
}

func (r *Repository) GetBookingByQRToken(ctx context.Context, qrToken string) (*models.Booking, error) {
	sql, args, err := psql.Select(bookingColumns...).From("bookings").
		Where(sq.Eq{"qr_token": qrToken}).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var b models.Booking
	if err := pgxscan.Get(ctx, r.db, &b, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "ticket not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get booking by qr token", err)
	}
	return &b, nil
}

func (r *Repository) MarkTicketUsed(ctx context.Context, bookingID string) error {
	sql, args, err := psql.Update("bookings").
		Set("used_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": bookingID, "used_at": nil}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	_, err = pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (string, error) {
		var id string
		return id, row.Scan(&id)
	})
	if err != nil {
		if pgxscan.NotFound(err) {
			return apperrors.New(apperrors.CodeAlreadyExists, "ticket already used")
		}
		return apperrors.Wrap(apperrors.CodeDBError, "mark ticket used", err)
	}
	return nil
}

func (r *Repository) FetchAndCancelExpiredBookings(ctx context.Context) ([]models.Booking, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "begin transaction", err)
	}
	defer tx.Rollback(ctx)

	updateSql, updateArgs, err := psql.Update("bookings").
		Set("status", utils.BookingStatusCancelled).
		Where(sq.And{
			sq.Eq{"status": utils.BookingStatusPending},
			sq.Expr("expires_at < NOW()"),
		}).
		Suffix("RETURNING " + pkgutils.JoinColumns(bookingColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build expire bookings query", err)
	}

	rows, _ := tx.Query(ctx, updateSql, updateArgs...)
	bookings, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Booking])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "cancel expired bookings", err)
	}

	if len(bookings) == 0 {
		return nil, nil
	}

	q := psql.Insert("bookings_state_log").Columns("booking_id", "from_status", "to_status")
	for _, b := range bookings {
		q = q.Values(b.ID, utils.BookingStatusPending, utils.BookingStatusCancelled)
	}
	logSql, logArgs, err := q.ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build state log query", err)
	}
	if _, err := tx.Exec(ctx, logSql, logArgs...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "insert state log entries", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "commit expire bookings", err)
	}

	return bookings, nil
}

func (r *Repository) ListBookingsByUserID(ctx context.Context, userID string) ([]models.Booking, error) {
	sql, args, err := psql.Select(bookingColumns...).From("bookings").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var bookings []models.Booking
	if err := pgxscan.Select(ctx, r.db, &bookings, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list bookings by user", err)
	}
	return bookings, nil
}
