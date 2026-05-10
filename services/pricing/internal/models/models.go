package models

import "time"

type Coupon struct {
	ID             string     `db:"id"`
	Code           string     `db:"code"`
	DiscountType   string     `db:"discount_type"`
	DiscountValue  string     `db:"discount_value"`
	MinOrderAmount string     `db:"min_order_amount"`
	MaxUses        *int       `db:"max_uses"`
	UsesRemaining  *int       `db:"uses_remaining"`
	MerchantID     string     `db:"merchant_id"`
	EventID        *string    `db:"event_id"`
	ValidFrom      time.Time  `db:"valid_from"`
	ValidUntil     time.Time  `db:"valid_until"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

type CouponRedemption struct {
	ID         string    `db:"id"`
	CouponID   string    `db:"coupon_id"`
	UserID     string    `db:"user_id"`
	BookingID  string    `db:"booking_id"`
	RedeemedAt time.Time `db:"redeemed_at"`
}

type PricingRule struct {
	ID         string    `db:"id"`
	MerchantID string    `db:"merchant_id"`
	RuleType   string    `db:"rule_type"`
	LowerBound *int      `db:"lower_bound"`
	UpperBound *int      `db:"upper_bound"`
	Multiplier string    `db:"multiplier"`
	Active     bool      `db:"active"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
