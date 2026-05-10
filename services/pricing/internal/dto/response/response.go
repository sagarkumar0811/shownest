package response

import "time"

type SeatPriceBreakdown struct {
	SeatID             string  `json:"seatId"`
	BasePrice          float64 `json:"basePrice"`
	CategoryMultiplier float64 `json:"categoryMultiplier"`
	SeatSubtotal       float64 `json:"seatSubtotal"`
}

type PriceCalculationResponse struct {
	ShowtimeID      string               `json:"showtimeId"`
	Seats           []SeatPriceBreakdown `json:"seats"`
	Subtotal        float64              `json:"subtotal"`
	TimeMultiplier  float64              `json:"timeMultiplier"`
	DayMultiplier   float64              `json:"dayMultiplier"`
	SurgeMultiplier float64              `json:"surgeMultiplier"`
	DynamicSubtotal float64              `json:"dynamicSubtotal"`
	CouponDiscount  float64              `json:"couponDiscount"`
	Total           float64              `json:"total"`
}

type CouponInfo struct {
	ID             string    `json:"id"`
	Code           string    `json:"code"`
	DiscountType   string    `json:"discountType"`
	DiscountValue  float64   `json:"discountValue"`
	MinOrderAmount float64   `json:"minOrderAmount"`
	MaxUses        *int      `json:"maxUses"`
	UsesRemaining  *int      `json:"usesRemaining"`
	EventID        *string   `json:"eventId,omitempty"`
	ValidFrom      time.Time `json:"validFrom"`
	ValidUntil     time.Time `json:"validUntil"`
}

type ValidateCouponResponse struct {
	Valid          bool        `json:"valid"`
	Coupon         *CouponInfo `json:"coupon,omitempty"`
	DiscountAmount float64     `json:"discountAmount"`
	Message        string      `json:"message,omitempty"`
}

type PricingRuleInfo struct {
	ID         string  `json:"id"`
	RuleType   string  `json:"ruleType"`
	LowerBound *int    `json:"lowerBound,omitempty"`
	UpperBound *int    `json:"upperBound,omitempty"`
	Multiplier float64 `json:"multiplier"`
	Active     bool    `json:"active"`
}
