package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	apperrors "github.com/shownest/pkg/errors"
	pkgutils "github.com/shownest/pkg/utils"
	"github.com/shownest/pricing-service/internal/models"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var couponColumns = []string{
	"id", "code", "discount_type", "discount_value", "min_order_amount",
	"max_uses", "uses_remaining", "merchant_id", "event_id",
	"valid_from", "valid_until", "created_at", "updated_at", "deleted_at",
}

var pricingRuleColumns = []string{
	"id", "merchant_id", "rule_type", "lower_bound", "upper_bound", "multiplier",
	"active", "created_at", "updated_at",
}

func (r *Repository) CreateCoupon(ctx context.Context, merchantID, code, discountType string, discountValue, minOrderAmount float64, maxUses *int, eventID *string, validFrom, validUntil time.Time) (*models.Coupon, error) {
	query, args, err := psql.
		Insert("coupons").
		Columns("code", "discount_type", "discount_value", "min_order_amount",
			"max_uses", "uses_remaining", "merchant_id", "event_id", "valid_from", "valid_until").
		Values(code, discountType, discountValue, minOrderAmount,
			maxUses, maxUses, merchantID, eventID, validFrom, validUntil).
		Suffix("RETURNING " + pkgutils.JoinColumns(couponColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build create coupon query", err)
	}

	var c models.Coupon
	if err := pgxscan.Get(ctx, r.db, &c, query, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create coupon", err)
	}
	return &c, nil
}

func (r *Repository) GetCouponByCode(ctx context.Context, code string) (*models.Coupon, error) {
	query, args, err := psql.
		Select(couponColumns...).
		From("coupons").
		Where(sq.Eq{"code": code, "deleted_at": nil}).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build get coupon query", err)
	}

	var c models.Coupon
	if err := pgxscan.Get(ctx, r.db, &c, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "coupon not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get coupon by code", err)
	}
	return &c, nil
}

func (r *Repository) DecrementCouponUses(ctx context.Context, couponID string) (*models.Coupon, error) {
	query, args, err := psql.
		Update("coupons").
		Set("uses_remaining", sq.Expr("uses_remaining - 1")).
		Where(sq.And{
			sq.Eq{"id": couponID, "deleted_at": nil},
			sq.Gt{"uses_remaining": 0},
		}).
		Suffix("RETURNING " + pkgutils.JoinColumns(couponColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build decrement coupon query", err)
	}

	var c models.Coupon
	if err := pgxscan.Get(ctx, r.db, &c, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeFailedPrecondition, "coupon has no uses remaining")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "decrement coupon uses", err)
	}
	return &c, nil
}

func (r *Repository) GetRedemptionByUser(ctx context.Context, couponID, userID string) (*models.CouponRedemption, error) {
	query, args, err := psql.
		Select("id", "coupon_id", "user_id", "booking_id", "redeemed_at").
		From("coupon_redemptions").
		Where(sq.Eq{"coupon_id": couponID, "user_id": userID}).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build get redemption query", err)
	}

	var red models.CouponRedemption
	if err := pgxscan.Get(ctx, r.db, &red, query, args...); err != nil {
		if pgxscan.NotFound(err) {
			return nil, apperrors.New(apperrors.CodeDBNotFound, "redemption not found")
		}
		return nil, apperrors.Wrap(apperrors.CodeDBError, "get coupon redemption", err)
	}
	return &red, nil
}

func (r *Repository) CreateRedemption(ctx context.Context, couponID, userID, bookingID string) (*models.CouponRedemption, error) {
	query, args, err := psql.
		Insert("coupon_redemptions").
		Columns("coupon_id", "user_id", "booking_id").
		Values(couponID, userID, bookingID).
		Suffix("RETURNING id, coupon_id, user_id, booking_id, redeemed_at").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build create redemption query", err)
	}

	var red models.CouponRedemption
	if err := pgxscan.Get(ctx, r.db, &red, query, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create coupon redemption", err)
	}
	return &red, nil
}

func (r *Repository) ListCouponsByMerchant(ctx context.Context, merchantID string) ([]models.Coupon, error) {
	query, args, err := psql.
		Select(couponColumns...).
		From("coupons").
		Where(sq.Eq{"merchant_id": merchantID, "deleted_at": nil}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build list coupons query", err)
	}

	var coupons []models.Coupon
	if err := pgxscan.Select(ctx, r.db, &coupons, query, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list coupons by merchant", err)
	}
	return coupons, nil
}

func (r *Repository) CreatePricingRule(ctx context.Context, merchantID, ruleType string, lowerBound, upperBound *int, multiplier float64) (*models.PricingRule, error) {
	query, args, err := psql.
		Insert("pricing_rules").
		Columns("merchant_id", "rule_type", "lower_bound", "upper_bound", "multiplier").
		Values(merchantID, ruleType, lowerBound, upperBound, multiplier).
		Suffix("RETURNING " + pkgutils.JoinColumns(pricingRuleColumns)).
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build create pricing rule query", err)
	}

	var rule models.PricingRule
	if err := pgxscan.Get(ctx, r.db, &rule, query, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "create pricing rule", err)
	}
	return &rule, nil
}

func (r *Repository) ListActivePricingRulesByMerchant(ctx context.Context, merchantID string) ([]models.PricingRule, error) {
	query, args, err := psql.
		Select(pricingRuleColumns...).
		From("pricing_rules").
		Where(sq.Eq{"merchant_id": merchantID, "active": true}).
		OrderBy("created_at").
		ToSql()
	if err != nil {
		return nil, apperrors.Wrap(apperrors.CodeInternal, "build list pricing rules query", err)
	}

	var rules []models.PricingRule
	if err := pgxscan.Select(ctx, r.db, &rules, query, args...); err != nil {
		return nil, apperrors.Wrap(apperrors.CodeDBError, "list active pricing rules", err)
	}
	return rules, nil
}
