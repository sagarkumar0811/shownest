-- Add pgcrypto extension for UUID generation and hashing
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Set the default time zone for the database
SET TIME ZONE 'Asia/Kolkata';

-- Users Table
CREATE TABLE IF NOT EXISTS users (
    id            UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    phone         VARCHAR(15) NOT NULL,
    email         VARCHAR(255),
    password_hash TEXT,
    google_id     VARCHAR(255),
    role          VARCHAR(20) NOT NULL DEFAULT 'Customer'
                  CHECK (role IN ('Customer', 'Merchant')),
    status        VARCHAR(20) NOT NULL DEFAULT 'Active'
                  CHECK (status IN ('Active', 'Blocked')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at    TIMESTAMPTZ,
    CONSTRAINT uq_users_phone     UNIQUE (phone),
    CONSTRAINT uq_users_email     UNIQUE (email),
    CONSTRAINT uq_users_google_id UNIQUE (google_id)
);

-- Indexes for Users Table
CREATE INDEX idx_users_phone     ON users (phone)     WHERE deleted_at IS NULL;
CREATE INDEX idx_users_email     ON users (email)     WHERE email IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_users_google_id ON users (google_id) WHERE google_id IS NOT NULL AND deleted_at IS NULL;
CREATE INDEX idx_users_status    ON users (status)    WHERE deleted_at IS NULL;

-- Trigger to update the updated_at timestamp on user updates
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for users table
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();

-- Sessions Table
CREATE TABLE IF NOT EXISTS sessions (
    id                 UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID        NOT NULL,
    refresh_token_hash CHAR(64)    NOT NULL,
    device_info        TEXT,
    ip_address         INET,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at         TIMESTAMPTZ NOT NULL,
    revoked_at         TIMESTAMPTZ
);

-- Indexes for Sessions Table
CREATE INDEX idx_sessions_user_id    ON sessions (user_id);
CREATE INDEX idx_sessions_token_hash ON sessions (refresh_token_hash);
CREATE INDEX idx_sessions_active     ON sessions (user_id, expires_at) WHERE revoked_at IS NULL;