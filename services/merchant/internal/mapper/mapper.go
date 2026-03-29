package mapper

import (
	"github.com/shownest/merchant-service/internal/dto/response"
	"github.com/shownest/merchant-service/internal/models"
)

func ToMerchantInfo(m *models.Merchant) response.MerchantInfo {
	return response.MerchantInfo{
		ID:           m.ID,
		UserID:       m.UserID,
		BusinessName: m.BusinessName,
		Category:     m.Category,
		ContactPhone: m.ContactPhone,
		ContactEmail: m.ContactEmail,
		Status:       m.Status,
		CreatedAt:    m.CreatedAt,
	}
}

func ToVenueInfo(v *models.Venue) response.VenueInfo {
	return response.VenueInfo{
		ID:         v.ID,
		MerchantID: v.MerchantID,
		Name:       v.Name,
		Address:    v.Address,
		City:       v.City,
		State:      v.State,
		Pincode:    v.Pincode,
		Latitude:   v.Latitude,
		Longitude:  v.Longitude,
		CreatedAt:  v.CreatedAt,
	}
}

func ToVenueInfoList(venues []models.Venue) []response.VenueInfo {
	infos := make([]response.VenueInfo, len(venues))
	for i, v := range venues {
		infos[i] = ToVenueInfo(&v)
	}
	return infos
}

func ToVenueWithDistanceInfo(v *models.VenueWithDistance) response.VenueWithDistanceInfo {
	return response.VenueWithDistanceInfo{
		VenueInfo:      ToVenueInfo(&v.Venue),
		DistanceMeters: v.DistanceMeters,
	}
}

func ToVenueWithDistanceInfoList(venues []models.VenueWithDistance) []response.VenueWithDistanceInfo {
	infos := make([]response.VenueWithDistanceInfo, len(venues))
	for i, v := range venues {
		infos[i] = ToVenueWithDistanceInfo(&v)
	}
	return infos
}

func ToHallInfo(h *models.Hall) response.HallInfo {
	return response.HallInfo{
		ID:        h.ID,
		VenueID:   h.VenueID,
		Name:      h.Name,
		Capacity:  h.Capacity,
		HallType:  h.HallType,
		CreatedAt: h.CreatedAt,
	}
}

func ToHallInfoList(halls []models.Hall) []response.HallInfo {
	infos := make([]response.HallInfo, len(halls))
	for i, h := range halls {
		infos[i] = ToHallInfo(&h)
	}
	return infos
}

func ToDocumentInfo(d *models.MerchantDocument) response.DocumentInfo {
	return response.DocumentInfo{
		ID:           d.ID,
		DocumentType: d.DocumentType,
		S3Key:        d.S3Key,
		VerifiedAt:   d.VerifiedAt,
		CreatedAt:    d.CreatedAt,
	}
}

func ToDocumentInfoList(docs []models.MerchantDocument) []response.DocumentInfo {
	infos := make([]response.DocumentInfo, len(docs))
	for i, d := range docs {
		infos[i] = ToDocumentInfo(&d)
	}
	return infos
}
