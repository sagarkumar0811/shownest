package request

import (
	"errors"
	"time"

	"github.com/shownest/pricing-service/internal/utils"
)

type CalculatePriceRequest struct {
	ShowtimeID string   `json:"showtimeId" binding:"required"`
	SeatIDs    []string `json:"seatIds" binding:"required,min=1"`
	CouponCode *string  `json:"couponCode"`
}

type ValidateCouponRequest struct {
	Code       string  `json:"code" binding:"required"`
	ShowtimeID string  `json:"showtimeId" binding:"required"`
	OrderTotal float64 `json:"orderTotal" binding:"required,gt=0"`
}

type CreateCouponRequest struct {
	Code           string    `json:"code" binding:"required"`
	DiscountType   string    `json:"discountType" binding:"required"`
	DiscountValue  float64   `json:"discountValue" binding:"required,gt=0"`
	MinOrderAmount float64   `json:"minOrderAmount"`
	MaxUses        *int      `json:"maxUses"`
	EventID        *string   `json:"eventId"`
	ValidFrom      time.Time `json:"validFrom" binding:"required"`
	ValidUntil     time.Time `json:"validUntil" binding:"required"`
}

func (r *CreateCouponRequest) Validate() error {
	if r.DiscountType != utils.DiscountTypeFlat && r.DiscountType != utils.DiscountTypePercent {
		return errors.New("discountType must be flat or percent")
	}
	if r.DiscountType == utils.DiscountTypePercent && r.DiscountValue > 100 {
		return errors.New("percent discount cannot exceed 100")
	}
	if r.ValidUntil.Before(r.ValidFrom) {
		return errors.New("validUntil must be after validFrom")
	}
	return nil
}

type CreatePricingRuleRequest struct {
	RuleType   string  `json:"ruleType" binding:"required"`
	LowerBound *int    `json:"lowerBound"`
	UpperBound *int    `json:"upperBound"`
	Multiplier float64 `json:"multiplier" binding:"required,gt=0"`
}

func (r *CreatePricingRuleRequest) Validate() error {
	switch r.RuleType {
	case utils.RuleTypeTimeOfDay:
		if r.LowerBound == nil || r.UpperBound == nil {
			return errors.New("hourly rule requires lowerBound (hour_from) and upperBound (hour_to)")
		}
		if *r.LowerBound < 0 || *r.LowerBound > 23 {
			return errors.New("lowerBound (hour_from) must be 0-23")
		}
		if *r.UpperBound < 1 || *r.UpperBound > 24 {
			return errors.New("upperBound (hour_to) must be 1-24")
		}
	case utils.RuleTypeDayOfWeek:
		if r.LowerBound == nil {
			return errors.New("weekly rule requires lowerBound (0=Sunday..6=Saturday)")
		}
		if *r.LowerBound < 0 || *r.LowerBound > 6 {
			return errors.New("lowerBound (day) must be 0-6")
		}
	case utils.RuleTypeSurge:
		if r.LowerBound == nil {
			return errors.New("surge rule requires lowerBound (occupancy threshold percent)")
		}
		if *r.LowerBound < 1 || *r.LowerBound > 100 {
			return errors.New("lowerBound (threshold) must be 1-100")
		}
	default:
		return errors.New("ruleType must be hourly, weekly, or surge")
	}
	return nil
}

type RecordRedemptionRequest struct {
	UserID    string `json:"userId" binding:"required"`
	BookingID string `json:"bookingId" binding:"required"`
}
