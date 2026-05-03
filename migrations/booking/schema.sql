CREATE EXTENSION IF NOT EXISTS "pgcrypto";

SET TIME ZONE 'Asia/Kolkata';

CREATE TABLE IF NOT EXISTS bookings (
    id           UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID           NOT NULL,
    showtime_id  UUID           NOT NULL,
    status       VARCHAR(20)    NOT NULL DEFAULT 'pending'
                 CHECK (status IN ('pending', 'confirmed', 'cancelled')),
    total_amount NUMERIC(10, 2) NOT NULL CHECK (total_amount >= 0),
    qr_token     TEXT,
    expires_at   TIMESTAMPTZ    NOT NULL,
    created_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id     ON bookings (user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_showtime_id ON bookings (showtime_id);
CREATE INDEX IF NOT EXISTS idx_bookings_status      ON bookings (status);

CREATE TABLE IF NOT EXISTS booking_items (
    id          UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id  UUID           NOT NULL REFERENCES bookings(id),
    seat_id     UUID           NOT NULL,
    category_id UUID           NOT NULL,
    price       NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_booking_items_booking_id ON booking_items (booking_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_booking_seat ON booking_items (booking_id, seat_id);

CREATE TABLE IF NOT EXISTS bookings_state_log (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id  UUID        NOT NULL REFERENCES bookings(id),
    from_status VARCHAR(20),
    to_status   VARCHAR(20) NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_state_log_booking_id ON bookings_state_log (booking_id);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_bookings_updated_at
    BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

GRANT ALL PRIVILEGES ON ALL TABLES    IN SCHEMA public TO booking_service;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO booking_service;
GRANT USAGE ON SCHEMA public TO booking_service;
