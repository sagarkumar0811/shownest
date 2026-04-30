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
	"id", "user_id", "showtime_id", "status", "total_amount", "qr_token", "expires_at", "created_at", "updated_at",
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
	q := psql.Update("bookings").
		Set("status", status).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING " + pkgutils.JoinColumns(bookingColumns))

	if qrToken != nil {
		q = q.Set("qr_token", *qrToken)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	b, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Booking])
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "booking not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "update booking status", err)
	}
	return &b, nil
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
