package models

import "time"

type Booking struct {
	ID          string     `db:"id"`
	UserID      string     `db:"user_id"`
	ShowtimeID  string     `db:"showtime_id"`
	Status      string     `db:"status"`
	TotalAmount float64    `db:"total_amount"`
	QRToken     *string    `db:"qr_token"`
	UsedAt      *time.Time `db:"used_at"`
	ExpiresAt   time.Time  `db:"expires_at"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

type BookingItem struct {
	ID         string    `db:"id"`
	BookingID  string    `db:"booking_id"`
	SeatID     string    `db:"seat_id"`
	CategoryID string    `db:"category_id"`
	Price      float64   `db:"price"`
	CreatedAt  time.Time `db:"created_at"`
}

type BookingStateLog struct {
	ID         string    `db:"id"`
	BookingID  string    `db:"booking_id"`
	FromStatus *string   `db:"from_status"`
	ToStatus   string    `db:"to_status"`
	CreatedAt  time.Time `db:"created_at"`
}
