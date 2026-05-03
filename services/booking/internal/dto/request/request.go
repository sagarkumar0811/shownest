package request

type BookingIDRequest struct {
	ID string `uri:"bookingId" binding:"required"`
}

type VerifyTicketRequest struct {
	Token string `json:"token" binding:"required"`
}

type SeatBookingItem struct {
	SeatID     string `json:"seatId"     binding:"required"`
	CategoryID string `json:"categoryId" binding:"required"`
}

type CreateBookingRequest struct {
	ShowtimeID string            `json:"showtimeId" binding:"required"`
	Seats      []SeatBookingItem `json:"seats"      binding:"required,min=1,dive"`
}
