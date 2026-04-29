package utils

const (
	SeatStatusAvailable = "available"
	SeatStatusLocked    = "locked"
	SeatStatusBooked    = "booked"

	SeatLockTTLSeconds = 600 // 10 minutes
	SeatAvailCacheTTL  = 10  // seconds
)

var ValidSeatStatuses = map[string]bool{
	SeatStatusAvailable: true,
	SeatStatusLocked:    true,
	SeatStatusBooked:    true,
}
