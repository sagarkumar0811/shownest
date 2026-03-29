package request

import (
	"errors"

	"github.com/shownest/merchant-service/internal/utils"
)

type CreateMerchantRequest struct {
	BusinessName string `json:"businessName" binding:"required,min=2,max=255"`
	Category     string `json:"category"     binding:"required"`
	ContactPhone string `json:"contactPhone" binding:"required"`
	ContactEmail string `json:"contactEmail" binding:"required,email"`
}

func (r *CreateMerchantRequest) Validate() error {
	if !utils.ValidCategories[r.Category] {
		return errors.New("invalid category")
	}
	return nil
}

type CreateVenueRequest struct {
	Name      string  `json:"name"      binding:"required,min=2,max=255"`
	Address   string  `json:"address"   binding:"required"`
	City      string  `json:"city"      binding:"required"`
	State     string  `json:"state"     binding:"required"`
	Pincode   string  `json:"pincode"   binding:"required"`
	Latitude  float64 `json:"latitude"  binding:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" binding:"required,min=-180,max=180"`
}

type CreateHallRequest struct {
	Name     string `json:"name"     binding:"required,min=2,max=255"`
	Capacity int    `json:"capacity" binding:"required,min=1"`
	HallType string `json:"hallType" binding:"required"`
}

func (r *CreateHallRequest) Validate() error {
	if !utils.ValidHallTypes[r.HallType] {
		return errors.New("invalid hall type")
	}
	return nil
}

type DocumentUploadURLRequest struct {
	DocumentType string `json:"documentType" binding:"required"`
}

func (r *DocumentUploadURLRequest) Validate() error {
	if !utils.ValidDocumentTypes[r.DocumentType] {
		return errors.New("invalid document type")
	}
	return nil
}

type ConfirmDocumentRequest struct {
	DocumentType string `json:"documentType" binding:"required"`
	S3Key        string `json:"s3Key"        binding:"required"`
}

func (r *ConfirmDocumentRequest) Validate() error {
	if !utils.ValidDocumentTypes[r.DocumentType] {
		return errors.New("invalid document type")
	}
	return nil
}

type NearbyVenuesRequest struct {
	Latitude  float64 `form:"latitude"  binding:"required,min=-90,max=90"`
	Longitude float64 `form:"longitude" binding:"required,min=-180,max=180"`
	Radius    float64 `form:"radius"    binding:"required,min=1"`
}

type VenueIDRequest struct {
	VenueID string `uri:"id" binding:"required"`
}

type CityRequest struct {
	City string `uri:"city" binding:"required"`
}
