package models

import "time"

type Merchant struct {
	ID           string     `db:"id"`
	UserID       string     `db:"user_id"`
	BusinessName string     `db:"business_name"`
	Category     string     `db:"category"`
	ContactPhone string     `db:"contact_phone"`
	ContactEmail string     `db:"contact_email"`
	Status       string     `db:"status"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type Venue struct {
	ID         string     `db:"id"`
	MerchantID string     `db:"merchant_id"`
	Name       string     `db:"name"`
	Address    string     `db:"address"`
	City       string     `db:"city"`
	State      string     `db:"state"`
	Pincode    string     `db:"pincode"`
	Latitude   float64    `db:"latitude"`
	Longitude  float64    `db:"longitude"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}

type VenueWithDistance struct {
	Venue
	DistanceMeters float64 `db:"distance_meters"`
}

type Hall struct {
	ID        string     `db:"id"`
	VenueID   string     `db:"venue_id"`
	Name      string     `db:"name"`
	Capacity  int        `db:"capacity"`
	HallType  string     `db:"hall_type"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type MerchantDocument struct {
	ID           string     `db:"id"`
	MerchantID   string     `db:"merchant_id"`
	DocumentType string     `db:"document_type"`
	S3Key        string     `db:"s3_key"`
	VerifiedAt   *time.Time `db:"verified_at"`
	CreatedAt    time.Time  `db:"created_at"`
}
