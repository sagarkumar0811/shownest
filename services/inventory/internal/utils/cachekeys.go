package utils

import "fmt"

const (
	cacheSeatAvailPrefix = "seat:avail:showtime:"
	lockSeatPrefix       = "seat:lock:"
)

// SeatAvailCacheKey is used to cache the full seat list for a showtime (10s TTL).
func SeatAvailCacheKey(showtimeID string) string {
	return cacheSeatAvailPrefix + showtimeID
}

// SeatLockKey is the Redis key used for atomic NX locking of a single seat.
func SeatLockKey(showtimeID, seatID string) string {
	return fmt.Sprintf("%s%s:%s", lockSeatPrefix, showtimeID, seatID)
}
