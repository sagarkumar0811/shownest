CREATE EXTENSION IF NOT EXISTS pgcrypto;

SET TIME ZONE 'Asia/Kolkata';

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE coupons (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code             TEXT NOT NULL UNIQUE,
    discount_type    TEXT NOT NULL CHECK (discount_type IN ('flat', 'percent')),
    discount_value   NUMERIC(10,2) NOT NULL CHECK (discount_value > 0),
    min_order_amount NUMERIC(10,2) NOT NULL DEFAULT 0,
    max_uses         INTEGER,
    uses_remaining   INTEGER,
    merchant_id      UUID         NOT NULL,
    event_id         UUID,
    valid_from       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    valid_until      TIMESTAMPTZ NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ
);

CREATE INDEX idx_coupons_code ON coupons(code) WHERE deleted_at IS NULL;
CREATE INDEX idx_coupons_merchant_id ON coupons(merchant_id) WHERE deleted_at IS NULL;

CREATE TABLE coupon_redemptions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    coupon_id   UUID NOT NULL REFERENCES coupons(id),
    user_id     UUID NOT NULL,
    booking_id  UUID NOT NULL UNIQUE,
    redeemed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_coupon_redemptions_unique ON coupon_redemptions(coupon_id, user_id);

-- Pricing rules can be of different types:
-- 1. Hourly: e.g. 10% off for bookings between 2pm-5pm (lower_bound=14, upper_bound=17)
-- 2. Weekly: e.g. 20% off for bookings on Saturday (lower_bound=6, upper_bound unused)
-- 3. Surge: e.g. 1.5x multiplier when occupancy exceeds a threshold (lower_bound=threshold%, upper_bound unused)
CREATE TABLE pricing_rules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL,
    rule_type   TEXT NOT NULL CHECK (rule_type IN ('hourly', 'weekly', 'surge')),
    lower_bound INTEGER,
    upper_bound INTEGER,
    multiplier  NUMERIC(5,3) NOT NULL CHECK (multiplier > 0),
    active      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pricing_rules_merchant_id ON pricing_rules (merchant_id) WHERE active = TRUE;

CREATE TRIGGER set_coupons_updated_at
    BEFORE UPDATE ON coupons
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER set_pricing_rules_updated_at
    BEFORE UPDATE ON pricing_rules
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
