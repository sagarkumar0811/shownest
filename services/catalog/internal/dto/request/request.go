package request

import (
	"errors"
	"time"

	"github.com/shownest/catalog-service/internal/utils"
)

type CreateEventRequest struct {
	MerchantID      string `json:"merchantId"   binding:"required"`
	Title           string `json:"title"        binding:"required,min=2,max=255"`
	Description     string `json:"description"`
	Category        string `json:"category"     binding:"required"`
	Language        string `json:"language"`
	DurationMinutes int    `json:"durationMinutes" binding:"required,min=1"`
	Rating          string `json:"rating"`
}

func (r *CreateEventRequest) Validate() error {
	if !utils.ValidCategories[r.Category] {
		return errors.New("invalid category")
	}
	return nil
}

type UpdateEventRequest struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	Language        *string `json:"language"`
	DurationMinutes *int    `json:"durationMinutes"`
	Rating          *string `json:"rating"`
	Status          *string `json:"status"`
}

func (r *UpdateEventRequest) Validate() error {
	if r.Status != nil && !utils.ValidEventStatuses[*r.Status] {
		return errors.New("invalid status")
	}
	if r.DurationMinutes != nil && *r.DurationMinutes < 1 {
		return errors.New("durationMinutes must be at least 1")
	}
	return nil
}

type EventIDRequest struct {
	EventID string `uri:"eventId" binding:"required"`
}

type ListEventsRequest struct {
	Category   string `form:"category"`
	MerchantID string `form:"merchantId"`
	Status     string `form:"status"`
}

type CreateShowtimeRequest struct {
	HallID    string    `json:"hallId"    binding:"required"`
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime"   binding:"required"`
	BasePrice float64   `json:"basePrice" binding:"required,min=0"`
}

func (r *CreateShowtimeRequest) Validate() error {
	if !r.EndTime.After(r.StartTime) {
		return errors.New("endTime must be after startTime")
	}
	if r.StartTime.Before(time.Now()) {
		return errors.New("startTime must be in the future")
	}
	return nil
}

type ShowtimeIDRequest struct {
	ShowtimeID string `uri:"showtimeId" binding:"required"`
}

type UpdateShowtimeRequest struct {
	Status *string `json:"status"`
}

func (r *UpdateShowtimeRequest) Validate() error {
	if r.Status != nil && !utils.ValidShowtimeStatuses[*r.Status] {
		return errors.New("invalid status")
	}
	return nil
}

type MediaUploadURLRequest struct {
	MediaType string `json:"mediaType" binding:"required"`
}

func (r *MediaUploadURLRequest) Validate() error {
	if !utils.ValidMediaTypes[r.MediaType] {
		return errors.New("invalid mediaType")
	}
	return nil
}

type ConfirmMediaRequest struct {
	MediaType string `json:"mediaType" binding:"required"`
	S3Key     string `json:"s3Key"     binding:"required"`
}

func (r *ConfirmMediaRequest) Validate() error {
	if !utils.ValidMediaTypes[r.MediaType] {
		return errors.New("invalid mediaType")
	}
	return nil
}
