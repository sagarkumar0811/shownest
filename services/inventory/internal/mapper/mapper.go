package mapper

import (
	"github.com/shownest/inventory-service/internal/dto/response"
	"github.com/shownest/inventory-service/internal/models"
)

func ToSeatCategoryInfo(c *models.SeatCategory) response.SeatCategoryInfo {
	return response.SeatCategoryInfo{
		ID:              c.ID,
		HallID:          c.HallID,
		Name:            c.Name,
		PriceMultiplier: c.PriceMultiplier,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
	}
}

func ToSeatCategoryInfoList(cats []models.SeatCategory) []response.SeatCategoryInfo {
	infos := make([]response.SeatCategoryInfo, len(cats))
	for i, c := range cats {
		infos[i] = ToSeatCategoryInfo(&c)
	}
	return infos
}

func ToSeatInfo(s *models.Seat) response.SeatInfo {
	return response.SeatInfo{
		ID:         s.ID,
		HallID:     s.HallID,
		CategoryID: s.CategoryID,
		Row:        s.Row,
		Number:     s.Number,
		XPosition:  s.XPosition,
		YPosition:  s.YPosition,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func ToSeatInfoList(seats []models.Seat) []response.SeatInfo {
	infos := make([]response.SeatInfo, len(seats))
	for i, s := range seats {
		infos[i] = ToSeatInfo(&s)
	}
	return infos
}

func ToShowtimeSeatInfo(ss *models.ShowtimeSeat, row string, number int, categoryID string) response.ShowtimeSeatInfo {
	return response.ShowtimeSeatInfo{
		ID:         ss.ID,
		ShowtimeID: ss.ShowtimeID,
		SeatID:     ss.SeatID,
		Row:        row,
		Number:     number,
		CategoryID: categoryID,
		Status:     ss.Status,
		LockedBy:   ss.LockedBy,
		LockedAt:   ss.LockedAt,
	}
}
