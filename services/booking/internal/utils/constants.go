package utils

const (
	BookingStatusPending   = "pending"
	BookingStatusConfirmed = "confirmed"
	BookingStatusCancelled = "cancelled"

	// BookingExpirySeconds defines how long a pending booking is valid
	// before it expires and the seat locks are released.
	BookingExpirySeconds = 600
)
