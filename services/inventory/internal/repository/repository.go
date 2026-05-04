package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shownest/inventory-service/internal/models"
	"github.com/shownest/inventory-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var seatCategoryColumns = []string{
	"id", "hall_id", "name", "price_multiplier", "created_at", "updated_at", "deleted_at",
}

var seatColumns = []string{
	"id", "hall_id", "category_id", "row", "number", "x_position", "y_position",
	"created_at", "updated_at", "deleted_at",
}

var showtimeSeatColumns = []string{
	"id", "showtime_id", "seat_id", "status", "locked_by", "locked_at", "created_at", "updated_at",
}

func (r *Repository) CreateSeatCategory(ctx context.Context, hallID, name string, multiplier float64) (*models.SeatCategory, error) {
	sql, args, err := psql.Insert("seat_categories").
		Columns("hall_id", "name", "price_multiplier").
		Values(hallID, name, multiplier).
		Suffix("RETURNING " + utils.JoinColumns(seatCategoryColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	c, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.SeatCategory])
	if err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeAlreadyExists, "seat category already exists for this hall")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create seat category", err)
	}
	return &c, nil
}

func (r *Repository) ListSeatCategoriesByHall(ctx context.Context, hallID string) ([]models.SeatCategory, error) {
	sql, args, err := psql.Select(seatCategoryColumns...).From("seat_categories").
		Where(sq.Eq{"hall_id": hallID}).
		Where("deleted_at IS NULL").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var cats []models.SeatCategory
	if err := pgxscan.Select(ctx, r.db, &cats, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list seat categories", err)
	}
	return cats, nil
}

func (r *Repository) GetSeatCategoryByID(ctx context.Context, id string) (*models.SeatCategory, error) {
	sql, args, err := psql.Select(seatCategoryColumns...).From("seat_categories").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var c models.SeatCategory
	if err := pgxscan.Get(ctx, r.db, &c, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "seat category not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get seat category", err)
	}
	return &c, nil
}

func (r *Repository) CreateSeat(ctx context.Context, hallID, categoryID, row string, number int, x, y *float64) (*models.Seat, error) {
	sql, args, err := psql.Insert("seats").
		Columns("hall_id", "category_id", "row", "number", "x_position", "y_position").
		Values(hallID, categoryID, row, number, nullFloat(x), nullFloat(y)).
		Suffix("RETURNING " + utils.JoinColumns(seatColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	s, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Seat])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create seat", err)
	}
	return &s, nil
}

func (r *Repository) BulkCreateSeats(ctx context.Context, hallID string, seats []struct {
	CategoryID string
	Row        string
	Number     int
	X, Y       *float64
}) ([]models.Seat, error) {
	q := psql.Insert("seats").Columns("hall_id", "category_id", "row", "number", "x_position", "y_position")
	for _, s := range seats {
		q = q.Values(hallID, s.CategoryID, s.Row, s.Number, nullFloat(s.X), nullFloat(s.Y))
	}
	q = q.Suffix("RETURNING " + utils.JoinColumns(seatColumns))

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	result, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Seat])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "bulk create seats", err)
	}
	return result, nil
}

func (r *Repository) ListSeatsByHall(ctx context.Context, hallID string) ([]models.Seat, error) {
	sql, args, err := psql.Select(seatColumns...).From("seats").
		Where(sq.Eq{"hall_id": hallID}).
		Where("deleted_at IS NULL").
		OrderBy("row ASC, number ASC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var result []models.Seat
	if err := pgxscan.Select(ctx, r.db, &result, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list seats by hall", err)
	}
	return result, nil
}

func (r *Repository) GetSeatByID(ctx context.Context, id string) (*models.Seat, error) {
	sql, args, err := psql.Select(seatColumns...).From("seats").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var s models.Seat
	if err := pgxscan.Get(ctx, r.db, &s, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "seat not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get seat", err)
	}
	return &s, nil
}

func (r *Repository) PublishShowtimeSeats(ctx context.Context, showtimeID, hallID string) (int, error) {
	seats, err := r.ListSeatsByHall(ctx, hallID)
	if err != nil {
		return 0, err
	}
	if len(seats) == 0 {
		return 0, apperrors.New(apperrors.CodeFailedPrecondition, "hall has no seats configured")
	}

	q := psql.Insert("showtime_seats").Columns("showtime_id", "seat_id").
		Suffix("ON CONFLICT (showtime_id, seat_id) DO NOTHING")
	for _, s := range seats {
		q = q.Values(showtimeID, s.ID)
	}

	sql, args, err := q.ToSql()
	if err != nil {
		return 0, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}
	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return 0, apperrors.Wrap(apperrors.CodeDBError, "publish showtime seats", err)
	}
	return len(seats), nil
}

func (r *Repository) ListShowtimeSeats(ctx context.Context, showtimeID string) ([]models.ShowtimeSeat, []models.Seat, error) {
	ssSql, ssArgs, err := psql.Select(showtimeSeatColumns...).From("showtime_seats").
		Where(sq.Eq{"showtime_id": showtimeID}).
		OrderBy("seat_id ASC").
		ToSql()
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var ssRows []models.ShowtimeSeat
	if err := pgxscan.Select(ctx, r.db, &ssRows, ssSql, ssArgs...); err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeDBError, "list showtime seats", err)
	}

	if len(ssRows) == 0 {
		return ssRows, nil, nil
	}

	seatIDs := make([]string, len(ssRows))
	for i, ss := range ssRows {
		seatIDs[i] = ss.SeatID
	}

	sSql, sArgs, err := psql.Select(seatColumns...).From("seats").
		Where(sq.Eq{"id": seatIDs}).
		ToSql()
	if err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var seatRows []models.Seat
	if err := pgxscan.Select(ctx, r.db, &seatRows, sSql, sArgs...); err != nil {
		return nil, nil, apperrors.Wrap(apperrors.CodeDBError, "fetch seat metadata", err)
	}

	return ssRows, seatRows, nil
}

func (r *Repository) LockShowtimeSeat(ctx context.Context, tx pgx.Tx, showtimeID, seatID, userID string) error {
	sql, args, err := psql.Update("showtime_seats").
		Set("status", utils.SeatStatusLocked).
		Set("locked_by", userID).
		Set("locked_at", sq.Expr("NOW()")).
		Where(sq.Eq{"showtime_id": showtimeID, "seat_id": seatID, "status": utils.SeatStatusAvailable}).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := tx.Query(ctx, sql, args...)
	_, err = pgx.CollectOneRow(rows, func(row pgx.CollectableRow) (string, error) {
		var id string
		return id, row.Scan(&id)
	})
	if err != nil {
		if pgxscan.NotFound(err) {
			return apperrors.New(apperrors.CodeAlreadyExists, "seat is no longer available")
		}
		return apperrors.Wrap(apperrors.CodeDBError, "lock showtime seat", err)
	}
	return nil
}

func (r *Repository) ReleaseShowtimeSeat(ctx context.Context, showtimeID, seatID, userID string) error {
	sql, args, err := psql.Update("showtime_seats").
		Set("status", utils.SeatStatusAvailable).
		Set("locked_by", nil).
		Set("locked_at", nil).
		Where(sq.Eq{"showtime_id": showtimeID, "seat_id": seatID, "status": utils.SeatStatusLocked, "locked_by": userID}).
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}
	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return apperrors.Wrap(apperrors.CodeDBError, "release showtime seat", err)
	}
	return nil
}

func (r *Repository) ExpireLockedSeats(ctx context.Context, ttlSeconds int) (int64, error) {
	sql, args, err := psql.Update("showtime_seats").
		Set("status", utils.SeatStatusAvailable).
		Set("locked_by", nil).
		Set("locked_at", nil).
		Where(sq.Eq{"status": utils.SeatStatusLocked}).
		Where(sq.Expr("locked_at < NOW() - INTERVAL '1 second' * ?", ttlSeconds)).
		ToSql()
	if err != nil {
		return 0, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return 0, apperrors.Wrap(apperrors.CodeDBError, "expire locked seats", err)
	}
	return tag.RowsAffected(), nil
}

func (r *Repository) BookShowtimeSeats(ctx context.Context, showtimeID, userID string, seatIDs []string) error {
	sql, args, err := psql.Update("showtime_seats").
		Set("status", utils.SeatStatusBooked).
		Set("locked_by", nil).
		Set("locked_at", nil).
		Where(sq.Eq{"showtime_id": showtimeID, "seat_id": seatIDs, "status": utils.SeatStatusLocked, "locked_by": userID}).
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}
	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return apperrors.Wrap(apperrors.CodeDBError, "book showtime seats", err)
	}
	return nil
}

func (r *Repository) GetSeatPrices(ctx context.Context, showtimeID string, seatIDs []string) ([]models.SeatPrice, error) {
	ifaces := make([]interface{}, len(seatIDs))
	for i, id := range seatIDs {
		ifaces[i] = id
	}

	sql, args, err := psql.
		Select("ss.seat_id", "s.category_id", "sc.price_multiplier").
		From("showtime_seats ss").
		Join("seats s ON s.id = ss.seat_id").
		Join("seat_categories sc ON sc.id = s.category_id").
		Where(sq.Eq{"ss.showtime_id": showtimeID}).
		Where(sq.Eq{"ss.seat_id": ifaces}).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var rows []models.SeatPrice
	if err := pgxscan.Select(ctx, r.db, &rows, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get seat prices", err)
	}
	return rows, nil
}

func (r *Repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "begin transaction", err)
	}
	return tx, nil
}

func (r *Repository) GetShowtimeOccupancy(ctx context.Context, showtimeID string) (*models.OccupancyStats, error) {
	query, args, err := psql.
		Select(
			"COUNT(*) AS total_seats",
			"COUNT(*) FILTER (WHERE status = 'booked') AS booked_seats",
		).
		From("showtime_seats").
		Where(sq.Eq{"showtime_id": showtimeID}).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build occupancy query", err)
	}

	var stats models.OccupancyStats
	if err := pgxscan.Get(ctx, r.db, &stats, query, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get showtime occupancy", err)
	}
	return &stats, nil
}

func nullFloat(v *float64) interface{} {
	if v == nil {
		return nil
	}
	return *v
}
