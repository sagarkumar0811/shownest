package utils

const (
	cacheEventPrefix    = "catalog:event:"
	cacheShowtimePrefix = "catalog:showtimes:event:"
)

func EventCacheKey(eventID string) string {
	return cacheEventPrefix + eventID
}

func ShowtimeListCacheKey(eventID string) string {
	return cacheShowtimePrefix + eventID
}
