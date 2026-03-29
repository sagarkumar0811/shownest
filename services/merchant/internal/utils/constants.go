package utils

const (
	MerchantStatusDraft     = "Draft"     // not submitted for review yet
	MerchantStatusPending   = "Pending"   // pending review
	MerchantStatusActive    = "Active"    // approved and active
	MerchantStatusRejected  = "Rejected"  // rejected after review
	MerchantStatusSuspended = "Suspended" // temporarily suspended
)

var ValidCategories = map[string]bool{
	"cinema":     true, // includes multiplexes and single-screen theaters
	"comedy":     true, // stand-up comedy venues
	"theatre":    true, // live performance theaters for plays, musicals, etc.
	"sports":     true, // venues for sporting events like football, basketball, etc.
	"music":      true, // concert halls, live music venues, etc.
	"dance":      true, // venues for dance performances
	"poetry":     true, // poetry reading venues
	"exhibition": true, // art and cultural exhibitions
	"other":      true, // other types of venues
}

var ValidHallTypes = map[string]bool{
	"auditorium":      true, // large indoor venue with tiered seating, suitable for concerts, theater, etc.
	"openstage":       true, // open stage for performances
	"lounge":          true, // lounge area for smaller gatherings or performances
	"outdoor":         true, // outdoor venue for events
	"arena":           true, // large arena for sports or concerts
	"multiplexscreen": true, // screen in a multiplex cinema
}

var ValidDocumentTypes = map[string]bool{
	"gstcertificate": true, // GST registration certificate
	"pan":            true, // PAN card
	"tradelicense":   true, // trade license
}

const DocumentUploadURLTTL = 15 // minutes
