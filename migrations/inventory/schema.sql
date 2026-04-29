CREATE EXTENSION IF NOT EXISTS "pgcrypto";

SET TIME ZONE 'Asia/Kolkata';

CREATE TABLE IF NOT EXISTS seat_categories (
    id               UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    hall_id          UUID           NOT NULL,
    name             VARCHAR(50)    NOT NULL,
    price_multiplier NUMERIC(5, 2)  NOT NULL DEFAULT 1.00 CHECK (price_multiplier > 0),
    created_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    deleted_at       TIMESTAMPTZ,
    CONSTRAINT uq_seat_category_hall_name UNIQUE (hall_id, name)
);

CREATE INDEX IF NOT EXISTS idx_seat_categories_hall_id ON seat_categories (hall_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS seats (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    hall_id     UUID        NOT NULL,
    category_id UUID        NOT NULL REFERENCES seat_categories(id),
    row         VARCHAR(10) NOT NULL,
    number      INT         NOT NULL CHECK (number > 0),
    x_position  NUMERIC(8, 2),
    y_position  NUMERIC(8, 2),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    CONSTRAINT uq_seat_hall_row_number UNIQUE (hall_id, row, number)
);

CREATE INDEX IF NOT EXISTS idx_seats_hall_id     ON seats (hall_id)     WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_seats_category_id ON seats (category_id) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS showtime_seats (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    showtime_id  UUID        NOT NULL,
    seat_id      UUID        NOT NULL REFERENCES seats(id),
    status       VARCHAR(20) NOT NULL DEFAULT 'available'
                 CHECK (status IN ('available', 'locked', 'booked')),
    locked_by    UUID,
    locked_at    TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_showtime_seat UNIQUE (showtime_id, seat_id)
);

CREATE INDEX IF NOT EXISTS idx_showtime_seats_showtime_id ON showtime_seats (showtime_id);
CREATE INDEX IF NOT EXISTS idx_showtime_seats_status      ON showtime_seats (showtime_id, status);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_seat_categories_updated_at
    BEFORE UPDATE ON seat_categories
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_seats_updated_at
    BEFORE UPDATE ON seats
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_showtime_seats_updated_at
    BEFORE UPDATE ON showtime_seats
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

GRANT USAGE ON SCHEMA public TO inventory_service;

GRANT SELECT, INSERT, UPDATE, DELETE
ON ALL TABLES IN SCHEMA public
TO inventory_service;

ALTER DEFAULT PRIVILEGES FOR ROLE postgres
IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE
ON TABLES TO inventory_service;
