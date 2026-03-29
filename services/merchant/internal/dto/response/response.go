package response

import "time"

type MerchantInfo struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	BusinessName string    `json:"businessName"`
	Category     string    `json:"category"`
	ContactPhone string    `json:"contactPhone"`
	ContactEmail string    `json:"contactEmail"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

type VenueInfo struct {
	ID         string    `json:"id"`
	MerchantID string    `json:"merchantId"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Pincode    string    `json:"pincode"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	CreatedAt  time.Time `json:"createdAt"`
}

type VenueWithDistanceInfo struct {
	VenueInfo
	DistanceMeters float64 `json:"distanceMeters"`
}

type HallInfo struct {
	ID        string    `json:"id"`
	VenueID   string    `json:"venueId"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	HallType  string    `json:"hallType"`
	CreatedAt time.Time `json:"createdAt"`
}

type DocumentInfo struct {
	ID           string     `json:"id"`
	DocumentType string     `json:"documentType"`
	S3Key        string     `json:"s3Key"`
	VerifiedAt   *time.Time `json:"verifiedAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type UploadURLResponse struct {
	UploadURL    string `json:"uploadUrl"`
	S3Key        string `json:"s3Key"`
	DocumentType string `json:"documentType"`
}
