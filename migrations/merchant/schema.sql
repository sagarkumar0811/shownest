CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "postgis";

SET TIME ZONE 'Asia/Kolkata';

CREATE TABLE IF NOT EXISTS merchants (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID         NOT NULL,
    business_name VARCHAR(255) NOT NULL,
    category      VARCHAR(50)  NOT NULL
                  CHECK (category IN ('cinema', 'comedy', 'theatre', 'sports', 'music', 'poetry', 'exhibition', 'other')),
    contact_phone VARCHAR(15)  NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    status        VARCHAR(30)  NOT NULL DEFAULT 'draft'
                  CHECK (status IN ('draft', 'pending_approval', 'active', 'rejected', 'suspended')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT uq_merchants_user_id UNIQUE (user_id)
);

CREATE INDEX idx_merchants_status   ON merchants (status)   WHERE deleted_at IS NULL;
CREATE INDEX idx_merchants_category ON merchants (category) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS venues (
    id          UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID             NOT NULL REFERENCES merchants(id),
    name        VARCHAR(255)     NOT NULL,
    address     TEXT             NOT NULL,
    city        VARCHAR(100)     NOT NULL,
    state       VARCHAR(100)     NOT NULL,
    pincode     VARCHAR(10)      NOT NULL,
    latitude    DOUBLE PRECISION NOT NULL,
    longitude   DOUBLE PRECISION NOT NULL,
    location    GEOGRAPHY(POINT, 4326) NOT NULL,
    created_at  TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX idx_venues_merchant_id ON venues (merchant_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_venues_city        ON venues (city)        WHERE deleted_at IS NULL;
CREATE INDEX idx_venues_location    ON venues USING GIST(location) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS halls (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    venue_id   UUID        NOT NULL REFERENCES venues(id),
    name       VARCHAR(255) NOT NULL,
    capacity   INTEGER      NOT NULL CHECK (capacity > 0),
    hall_type  VARCHAR(30)  NOT NULL
               CHECK (hall_type IN ('auditorium', 'open_stage', 'lounge', 'outdoor', 'arena', 'multiplex_screen')),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_halls_venue_id ON halls (venue_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS merchant_documents (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id   UUID        NOT NULL REFERENCES merchants(id),
    document_type VARCHAR(50) NOT NULL
                  CHECK (document_type IN ('gst_certificate', 'pan', 'trade_license')),
    s3_key        TEXT        NOT NULL,
    verified_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_merchant_documents_merchant_id ON merchant_documents (merchant_id);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_merchants_updated_at
    BEFORE UPDATE ON merchants
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_venues_updated_at
    BEFORE UPDATE ON venues
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_halls_updated_at
    BEFORE UPDATE ON halls
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

GRANT USAGE ON SCHEMA public TO merchant_service;

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA public
TO merchant_service;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres
IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE
ON TABLES TO merchant_service;
