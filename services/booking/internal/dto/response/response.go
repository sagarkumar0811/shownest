package response

import "time"

type BookingItemInfo struct {
	ID         string    `json:"id"`
	SeatID     string    `json:"seatId"`
	CategoryID string    `json:"categoryId"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"createdAt"`
}

type BookingInfo struct {
	ID          string            `json:"id"`
	UserID      string            `json:"userId"`
	ShowtimeID  string            `json:"showtimeId"`
	Status      string            `json:"status"`
	TotalAmount float64           `json:"totalAmount"`
	QRToken     *string           `json:"qrToken,omitempty"`
	ExpiresAt   time.Time         `json:"expiresAt"`
	Items       []BookingItemInfo `json:"items"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}
