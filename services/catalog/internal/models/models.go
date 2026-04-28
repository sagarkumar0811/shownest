package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Event struct {
	ID              string     `db:"id"`
	UserID          string     `db:"user_id"`
	MerchantID      string     `db:"merchant_id"`
	Title           string     `db:"title"`
	Description     *string    `db:"description"`
	Category        string     `db:"category"`
	Language        *string    `db:"language"`
	DurationMinutes int        `db:"duration_minutes"`
	Rating          *string    `db:"rating"`
	Status          string     `db:"status"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

type Showtime struct {
	ID        string         `db:"id"`
	EventID   string         `db:"event_id"`
	HallID    string         `db:"hall_id"`
	StartTime time.Time      `db:"start_time"`
	EndTime   time.Time      `db:"end_time"`
	BasePrice pgtype.Numeric `db:"base_price"`
	Status    string         `db:"status"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	DeletedAt *time.Time     `db:"deleted_at"`
}

type EventMedia struct {
	ID        string    `db:"id"`
	EventID   string    `db:"event_id"`
	MediaType string    `db:"media_type"`
	S3Key     string    `db:"s3_key"`
	CdnURL    *string   `db:"cdn_url"`
	CreatedAt time.Time `db:"created_at"`
}
