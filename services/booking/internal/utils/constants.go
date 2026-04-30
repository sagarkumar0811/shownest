package utils

const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCancelled = "cancelled"
	BookingStatusExpired   = "expired"

	// BookingExpirySeconds defines how long a pending booking is valid
	// before it expires and the seat locks are released.
	BookingExpirySeconds = 600
)
