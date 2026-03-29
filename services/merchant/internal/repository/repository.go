package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shownest/merchant-service/internal/models"
	"github.com/shownest/merchant-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var merchantColumns = []string{
	"id", "user_id", "business_name", "category", "contact_phone", "contact_email",
	"status", "created_at", "updated_at", "deleted_at",
}

var venueColumns = []string{
	"id", "merchant_id", "name", "address", "city", "state", "pincode",
	"latitude", "longitude", "created_at", "updated_at", "deleted_at",
}

var hallColumns = []string{
	"id", "venue_id", "name", "capacity", "hall_type",
	"created_at", "updated_at", "deleted_at",
}

var documentColumns = []string{
	"id", "merchant_id", "document_type", "s3_key", "verified_at", "created_at",
}

func (r *Repository) CreateMerchant(ctx context.Context, userID, businessName, category, contactPhone, contactEmail string) (*models.Merchant, error) {
	sql, args, err := psql.Insert("merchants").
		Columns("user_id", "business_name", "category", "contact_phone", "contact_email").
		Values(userID, businessName, category, contactPhone, contactEmail).
		Suffix("RETURNING " + utils.JoinColumns(merchantColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	m, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Merchant])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create merchant", err)
	}
	return &m, nil
}

func (r *Repository) GetMerchantByUserID(ctx context.Context, userID string) (*models.Merchant, error) {
	sql, args, err := psql.Select(merchantColumns...).From("merchants").
		Where(sq.Eq{"user_id": userID}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var m models.Merchant
	if err := pgxscan.Get(ctx, r.db, &m, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "merchant not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get merchant by user id", err)
	}
	return &m, nil
}

func (r *Repository) GetMerchantByID(ctx context.Context, id string) (*models.Merchant, error) {
	sql, args, err := psql.Select(merchantColumns...).From("merchants").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var m models.Merchant
	if err := pgxscan.Get(ctx, r.db, &m, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "merchant not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get merchant by id", err)
	}
	return &m, nil
}

func (r *Repository) UpdateMerchantStatus(ctx context.Context, id, status string) error {
	sql, args, err := psql.Update("merchants").
		Set("status", status).
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeDBError, "update merchant status", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.CodeDBNotFound, "merchant not found")
	}
	return nil
}

func (r *Repository) CreateVenue(ctx context.Context, merchantID, name, address, city, state, pincode string, lat, lng float64) (*models.Venue, error) {
	rawSQL := fmt.Sprintf(`
		INSERT INTO venues (merchant_id, name, address, city, state, pincode, latitude, longitude, location)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, ST_MakePoint($9, $10)::geography)
		RETURNING %s`, utils.JoinColumns(venueColumns))

	rows, _ := r.db.Query(ctx, rawSQL, merchantID, name, address, city, state, pincode, lat, lng, lng, lat)
	v, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Venue])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create venue", err)
	}
	return &v, nil
}

func (r *Repository) GetVenueByID(ctx context.Context, id string) (*models.Venue, error) {
	sql, args, err := psql.Select(venueColumns...).From("venues").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var v models.Venue
	if err := pgxscan.Get(ctx, r.db, &v, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "venue not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get venue by id", err)
	}
	return &v, nil
}

func (r *Repository) GetVenuesByMerchantID(ctx context.Context, merchantID string) ([]models.Venue, error) {
	sql, args, err := psql.Select(venueColumns...).From("venues").
		Where(sq.Eq{"merchant_id": merchantID}).
		Where("deleted_at IS NULL").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var venues []models.Venue
	if err := pgxscan.Select(ctx, r.db, &venues, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list venues by merchant", err)
	}
	return venues, nil
}

func (r *Repository) GetNearbyVenues(ctx context.Context, lat, lng, radiusMeters float64) ([]models.VenueWithDistance, error) {
	rawSQL := fmt.Sprintf(`
		SELECT %s, ST_Distance(location, ST_MakePoint($2, $1)::geography) AS distance_meters
		FROM venues
		WHERE ST_DWithin(location, ST_MakePoint($2, $1)::geography, $3)
		  AND deleted_at IS NULL
		ORDER BY distance_meters`, utils.JoinColumns(venueColumns))

	var venues []models.VenueWithDistance
	if err := pgxscan.Select(ctx, r.db, &venues, rawSQL, lat, lng, radiusMeters); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get nearby venues", err)
	}
	return venues, nil
}

func (r *Repository) GetVenuesByCity(ctx context.Context, city string) ([]models.Venue, error) {
	sql, args, err := psql.Select(venueColumns...).From("venues").
		Where(sq.ILike{"city": city}).
		Where("deleted_at IS NULL").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var venues []models.Venue
	if err := pgxscan.Select(ctx, r.db, &venues, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get venues by city", err)
	}
	return venues, nil
}

func (r *Repository) CreateHall(ctx context.Context, venueID, name string, capacity int, hallType string) (*models.Hall, error) {
	sql, args, err := psql.Insert("halls").
		Columns("venue_id", "name", "capacity", "hall_type").
		Values(venueID, name, capacity, hallType).
		Suffix("RETURNING " + utils.JoinColumns(hallColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	h, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Hall])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create hall", err)
	}
	return &h, nil
}

func (r *Repository) GetHallsByVenueID(ctx context.Context, venueID string) ([]models.Hall, error) {
	sql, args, err := psql.Select(hallColumns...).From("halls").
		Where(sq.Eq{"venue_id": venueID}).
		Where("deleted_at IS NULL").
		OrderBy("name ASC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var halls []models.Hall
	if err := pgxscan.Select(ctx, r.db, &halls, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list halls by venue", err)
	}
	return halls, nil
}

func (r *Repository) CreateDocument(ctx context.Context, merchantID, docType, s3Key string) (*models.MerchantDocument, error) {
	sql, args, err := psql.Insert("merchant_documents").
		Columns("merchant_id", "document_type", "s3_key").
		Values(merchantID, docType, s3Key).
		Suffix("RETURNING " + utils.JoinColumns(documentColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	d, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.MerchantDocument])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create document", err)
	}
	return &d, nil
}

func (r *Repository) GetDocumentsByMerchantID(ctx context.Context, merchantID string) ([]models.MerchantDocument, error) {
	sql, args, err := psql.Select(documentColumns...).From("merchant_documents").
		Where(sq.Eq{"merchant_id": merchantID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var docs []models.MerchantDocument
	if err := pgxscan.Select(ctx, r.db, &docs, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list documents by merchant", err)
	}
	return docs, nil
}
