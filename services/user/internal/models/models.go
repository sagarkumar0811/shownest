package models

import "time"

type User struct {
	ID           string     `db:"id"`
	Phone        string     `db:"phone"`
	Email        *string    `db:"email"`
	PasswordHash *string    `db:"password_hash"`
	GoogleID     *string    `db:"google_id"`
	Role         string     `db:"role"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type Session struct {
	ID               string     `db:"id"`
	UserID           string     `db:"user_id"`
	RefreshTokenHash string     `db:"refresh_token_hash"`
	DeviceInfo       *string    `db:"device_info"`
	IPAddress        *string    `db:"ip_address"`
	CreatedAt        time.Time  `db:"created_at"`
	ExpiresAt        time.Time  `db:"expires_at"`
	RevokedAt        *time.Time `db:"revoked_at"`
}
