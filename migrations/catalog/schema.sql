CREATE EXTENSION IF NOT EXISTS "pgcrypto";

SET TIME ZONE 'Asia/Kolkata';

CREATE TABLE IF NOT EXISTS events (
    id            UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID         NOT NULL,
    merchant_id   UUID         NOT NULL,
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    category      VARCHAR(50)  NOT NULL
                  CHECK (category IN ('cinema', 'comedy', 'theatre', 'sports', 'music', 'dance', 'poetry', 'exhibition', 'other')),
    language      VARCHAR(100),
    duration_minutes INTEGER      NOT NULL CHECK (duration_minutes > 0),
    rating        VARCHAR(10),
    status        VARCHAR(30)  NOT NULL DEFAULT 'draft'
                  CHECK (status IN ('draft', 'published', 'cancelled')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_events_merchant_id ON events (merchant_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_events_category    ON events (category)    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_events_status      ON events (status)      WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS showtimes (
    id           UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id     UUID           NOT NULL REFERENCES events(id),
    hall_id      UUID           NOT NULL,
    start_time   TIMESTAMPTZ    NOT NULL,
    end_time     TIMESTAMPTZ    NOT NULL,
    base_price   NUMERIC(10, 2) NOT NULL CHECK (base_price >= 0),
    status       VARCHAR(30)    NOT NULL DEFAULT 'scheduled'
                 CHECK (status IN ('scheduled', 'cancelled', 'completed')),
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at   TIMESTAMPTZ,
    CONSTRAINT chk_showtime_times CHECK (end_time > start_time)
);

CREATE INDEX IF NOT EXISTS idx_showtimes_event_id   ON showtimes (event_id)   WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_showtimes_hall_id    ON showtimes (hall_id)    WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_showtimes_start_time ON showtimes (start_time) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS event_media (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id   UUID        NOT NULL REFERENCES events(id),
    media_type VARCHAR(20) NOT NULL
               CHECK (media_type IN ('poster', 'trailer')),
    s3_key     TEXT        NOT NULL,
    cdn_url    TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_event_media_event_id ON event_media (event_id);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_events_updated_at
    BEFORE UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_showtimes_updated_at
    BEFORE UPDATE ON showtimes
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
