package mapper

import (
	"github.com/shownest/booking-service/internal/dto/response"
	"github.com/shownest/booking-service/internal/models"
)

func ToBookingItemInfo(item *models.BookingItem) response.BookingItemInfo {
	return response.BookingItemInfo{
		ID:         item.ID,
		SeatID:     item.SeatID,
		CategoryID: item.CategoryID,
		Price:      item.Price,
		CreatedAt:  item.CreatedAt,
	}
}

func ToBookingInfo(b *models.Booking, items []models.BookingItem) response.BookingInfo {
	itemInfos := make([]response.BookingItemInfo, len(items))
	for i, item := range items {
		itemInfos[i] = ToBookingItemInfo(&item)
	}
	return response.BookingInfo{
		ID:          b.ID,
		UserID:      b.UserID,
		ShowtimeID:  b.ShowtimeID,
		Status:      b.Status,
		TotalAmount: b.TotalAmount,
		QRToken:     b.QRToken,
		ExpiresAt:   b.ExpiresAt,
		Items:       itemInfos,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func ToBookingInfoList(bookings []models.Booking) []response.BookingInfo {
	infos := make([]response.BookingInfo, len(bookings))
	for i, b := range bookings {
		infos[i] = ToBookingInfo(&b, nil)
	}
	return infos
}
