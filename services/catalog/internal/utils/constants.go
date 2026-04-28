package utils

const (
	EventStatusDraft     = "draft"
	EventStatusPublished = "published"
	EventStatusCancelled = "cancelled"

	ShowtimeStatusScheduled = "scheduled"
	ShowtimeStatusCancelled = "cancelled"
	ShowtimeStatusCompleted = "completed"

	MediaUploadURLTTL = 15
)

var ValidCategories = map[string]bool{
	"cinema":     true,
	"comedy":     true,
	"theatre":    true,
	"sports":     true,
	"music":      true,
	"dance":      true,
	"poetry":     true,
	"exhibition": true,
	"other":      true,
}

var ValidEventStatuses = map[string]bool{
	EventStatusDraft:     true,
	EventStatusPublished: true,
	EventStatusCancelled: true,
}

var ValidShowtimeStatuses = map[string]bool{
	ShowtimeStatusScheduled: true,
	ShowtimeStatusCancelled: true,
	ShowtimeStatusCompleted: true,
}

var ValidMediaTypes = map[string]bool{
	"poster":  true,
	"trailer": true,
}
