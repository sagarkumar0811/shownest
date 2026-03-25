package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/user-service/internal/models"
	"github.com/shownest/user-service/internal/utils"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var userColumns = []string{
	"id", "phone", "email", "password_hash", "google_id",
	"role", "status", "created_at", "updated_at", "deleted_at",
}

var sessionColumns = []string{
	"id", "user_id", "refresh_token_hash", "device_info", "ip_address",
	"created_at", "expires_at", "revoked_at",
}

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	sql, args, err := psql.Select(userColumns...).From("users").
		Where(sq.Eq{"phone": phone}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var u models.User
	if err := pgxscan.Get(ctx, r.db, &u, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "user not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get user by phone", err)
	}
	return &u, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	sql, args, err := psql.Select(userColumns...).From("users").
		Where(sq.Eq{"id": id}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var u models.User
	if err := pgxscan.Get(ctx, r.db, &u, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "user not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get user by id", err)
	}
	return &u, nil
}

func (r *Repository) CreateUser(ctx context.Context, phone, role string) (*models.User, error) {
	sql, args, err := psql.Insert("users").
		Columns("phone", "role").
		Values(phone, role).
		Suffix("RETURNING " + utils.JoinColumns(userColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	u, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.User])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create user", err)
	}
	return &u, nil
}

func (r *Repository) UpdateUserEmail(ctx context.Context, userID string, email *string) error {
	sql, args, err := psql.Update("users").
		Set("email", email).
		Where(sq.Eq{"id": userID}).
		Where("deleted_at IS NULL").
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return apperrors.Wrap(apperrors.CodeDBError, "update user email", err)
	}
	return nil
}

func (r *Repository) CreateSession(ctx context.Context, id, userID, tokenHash, deviceInfo, ipAddress string, expiresAt time.Time) (*models.Session, error) {
	sql, args, err := psql.Insert("sessions").
		Columns("id", "user_id", "refresh_token_hash", "device_info", "ip_address", "expires_at").
		Values(id, userID, tokenHash, deviceInfo, ipAddress, expiresAt).
		Suffix("RETURNING " + utils.JoinColumns(sessionColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	rows, _ := r.db.Query(ctx, sql, args...)
	s, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Session])
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create session", err)
	}
	return &s, nil
}

func (r *Repository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (*models.Session, error) {
	sql, args, err := psql.Select(sessionColumns...).From("sessions").
		Where(sq.Eq{"refresh_token_hash": tokenHash}).
		Where("revoked_at IS NULL").
		Where("expires_at > NOW()").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var s models.Session
	if err := pgxscan.Get(ctx, r.db, &s, sql, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "session not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get session by token hash", err)
	}
	return &s, nil
}

func (r *Repository) ListActiveSessions(ctx context.Context, userID string) ([]models.Session, error) {
	sql, args, err := psql.Select(sessionColumns...).From("sessions").
		Where(sq.Eq{"user_id": userID}).
		Where("revoked_at IS NULL").
		Where("expires_at > NOW()").
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	var sessions []models.Session
	if err := pgxscan.Select(ctx, r.db, &sessions, sql, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list active sessions", err)
	}
	return sessions, nil
}

func (r *Repository) RevokeSession(ctx context.Context, sessionID, userID string) error {
	sql, args, err := psql.Update("sessions").
		Set("revoked_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": sessionID, "user_id": userID}).
		Where("revoked_at IS NULL").
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeDBError, "revoke session", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.CodeDBNotFound, "session not found")
	}
	return nil
}

func (r *Repository) RevokeAllSessionsExcept(ctx context.Context, userID, currentSessionID string) error {
	sql, args, err := psql.Update("sessions").
		Set("revoked_at", sq.Expr("NOW()")).
		Where(sq.Eq{"user_id": userID}).
		Where(sq.NotEq{"id": currentSessionID}).
		Where("revoked_at IS NULL").
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	if _, err := r.db.Exec(ctx, sql, args...); err != nil {
		return apperrors.Wrap(apperrors.CodeDBError, "revoke all sessions", err)
	}
	return nil
}

func (r *Repository) RotateSessionToken(ctx context.Context, sessionID, newTokenHash string, newExpiresAt time.Time) error {
	sql, args, err := psql.Update("sessions").
		Set("refresh_token_hash", newTokenHash).
		Set("expires_at", newExpiresAt).
		Where(sq.Eq{"id": sessionID}).
		Where("revoked_at IS NULL").
		ToSql()
	if err != nil {
		return apperrors.Wrap(apperrors.CodeInternal, "build query", err)
	}

	tag, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return apperrors.Wrap(apperrors.CodeDBError, "rotate session token", err)
	}
	if tag.RowsAffected() == 0 {
		return apperrors.New(apperrors.CodeDBNotFound, "session not found or already revoked")
	}
	return nil
}
