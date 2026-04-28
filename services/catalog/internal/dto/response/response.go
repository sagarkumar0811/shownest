package response

import "time"

type EventInfo struct {
	ID              string    `json:"id"`
	UserID          string    `json:"userId"`
	MerchantID      string    `json:"merchantId"`
	Title           string    `json:"title"`
	Description     *string   `json:"description,omitempty"`
	Category        string    `json:"category"`
	Language        *string   `json:"language,omitempty"`
	DurationMinutes int       `json:"durationMinutes"`
	Rating          *string   `json:"rating,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type ShowtimeInfo struct {
	ID        string    `json:"id"`
	EventID   string    `json:"eventId"`
	HallID    string    `json:"hallId"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	BasePrice string    `json:"basePrice"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MediaInfo struct {
	ID        string    `json:"id"`
	EventID   string    `json:"eventId"`
	MediaType string    `json:"mediaType"`
	S3Key     string    `json:"s3Key"`
	CdnURL    *string   `json:"cdnUrl,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type MediaUploadURLResponse struct {
	UploadURL string `json:"uploadUrl"`
	S3Key     string `json:"s3Key"`
	MediaType string `json:"mediaType"`
}
