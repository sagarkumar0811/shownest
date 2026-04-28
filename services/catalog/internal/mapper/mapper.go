package mapper

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shownest/catalog-service/internal/dto/response"
	"github.com/shownest/catalog-service/internal/models"
)

func ToEventInfo(e *models.Event) response.EventInfo {
	return response.EventInfo{
		ID:              e.ID,
		UserID:          e.UserID,
		MerchantID:      e.MerchantID,
		Title:           e.Title,
		Description:     e.Description,
		Category:        e.Category,
		Language:        e.Language,
		DurationMinutes: e.DurationMinutes,
		Rating:          e.Rating,
		Status:          e.Status,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}

func ToEventInfoList(events []models.Event) []response.EventInfo {
	infos := make([]response.EventInfo, len(events))
	for i, e := range events {
		infos[i] = ToEventInfo(&e)
	}
	return infos
}

func ToShowtimeInfo(s *models.Showtime) response.ShowtimeInfo {
	return response.ShowtimeInfo{
		ID:        s.ID,
		EventID:   s.EventID,
		HallID:    s.HallID,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		BasePrice: numericToString(s.BasePrice),
		Status:    s.Status,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func ToShowtimeInfoList(showtimes []models.Showtime) []response.ShowtimeInfo {
	infos := make([]response.ShowtimeInfo, len(showtimes))
	for i, s := range showtimes {
		infos[i] = ToShowtimeInfo(&s)
	}
	return infos
}

func ToMediaInfo(m *models.EventMedia) response.MediaInfo {
	return response.MediaInfo{
		ID:        m.ID,
		EventID:   m.EventID,
		MediaType: m.MediaType,
		S3Key:     m.S3Key,
		CdnURL:    m.CdnURL,
		CreatedAt: m.CreatedAt,
	}
}

func ToMediaInfoList(media []models.EventMedia) []response.MediaInfo {
	infos := make([]response.MediaInfo, len(media))
	for i, m := range media {
		infos[i] = ToMediaInfo(&m)
	}
	return infos
}

func numericToString(n pgtype.Numeric) string {
	if !n.Valid {
		return "0"
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return "0"
	}
	return fmt.Sprintf("%.2f", f.Float64)
}
