package models

import "time"

type SeatCategory struct {
	ID              string     `db:"id"`
	HallID          string     `db:"hall_id"`
	Name            string     `db:"name"`
	PriceMultiplier string     `db:"price_multiplier"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
}

type Seat struct {
	ID         string     `db:"id"`
	HallID     string     `db:"hall_id"`
	CategoryID string     `db:"category_id"`
	Row        string     `db:"row"`
	Number     int        `db:"number"`
	XPosition  *string    `db:"x_position"`
	YPosition  *string    `db:"y_position"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

type ShowtimeSeat struct {
	ID         string     `db:"id"`
	ShowtimeID string     `db:"showtime_id"`
	SeatID     string     `db:"seat_id"`
	Status     string     `db:"status"`
	LockedBy   *string    `db:"locked_by"`
	LockedAt   *time.Time `db:"locked_at"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}

type SeatPrice struct {
	SeatID          string `db:"seat_id"`
	CategoryID      string `db:"category_id"`
	PriceMultiplier string `db:"price_multiplier"`
}
