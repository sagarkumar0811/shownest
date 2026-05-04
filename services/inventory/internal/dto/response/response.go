package response

import "time"

type SeatCategoryInfo struct {
	ID              string    `json:"id"`
	HallID          string    `json:"hallId"`
	Name            string    `json:"name"`
	PriceMultiplier string    `json:"priceMultiplier"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type SeatInfo struct {
	ID         string    `json:"id"`
	HallID     string    `json:"hallId"`
	CategoryID string    `json:"categoryId"`
	Row        string    `json:"row"`
	Number     int       `json:"number"`
	XPosition  *string   `json:"xPosition,omitempty"`
	YPosition  *string   `json:"yPosition,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type ShowtimeSeatInfo struct {
	ID         string     `json:"id"`
	ShowtimeID string     `json:"showtimeId"`
	SeatID     string     `json:"seatId"`
	Row        string     `json:"row"`
	Number     int        `json:"number"`
	CategoryID string     `json:"categoryId"`
	Status     string     `json:"status"`
	LockedBy   *string    `json:"lockedBy,omitempty"`
	LockedAt   *time.Time `json:"lockedAt,omitempty"`
}

type LockSeatsResponse struct {
	Locked []string `json:"locked"`
	Failed []string `json:"failed,omitempty"`
}

type SeatPriceInfo struct {
	SeatID          string `json:"seatId"`
	CategoryID      string `json:"categoryId"`
	PriceMultiplier string `json:"priceMultiplier"`
}

type OccupancyInfo struct {
	TotalSeats       int     `json:"totalSeats"`
	BookedSeats      int     `json:"bookedSeats"`
	OccupancyPercent float64 `json:"occupancyPercent"`
}
