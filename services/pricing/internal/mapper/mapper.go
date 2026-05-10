package mapper

import (
	"strconv"

	"github.com/shownest/pricing-service/internal/client"
	"github.com/shownest/pricing-service/internal/dto/response"
	"github.com/shownest/pricing-service/internal/models"
)

func ToSeatPriceBreakdownList(basePrice float64, seatPrices []client.SeatPriceInfo) ([]response.SeatPriceBreakdown, float64) {
	breakdowns := make([]response.SeatPriceBreakdown, 0, len(seatPrices))
	var subtotal float64
	for _, sp := range seatPrices {
		priceMultiplier, _ := strconv.ParseFloat(sp.PriceMultiplier, 64)
		seatSubtotal := basePrice * priceMultiplier
		subtotal += seatSubtotal
		breakdowns = append(breakdowns, response.SeatPriceBreakdown{
			SeatID:             sp.SeatID,
			BasePrice:          basePrice,
			CategoryMultiplier: priceMultiplier,
			SeatSubtotal:       seatSubtotal,
		})
	}
	return breakdowns, subtotal
}

func ToCouponInfo(c *models.Coupon) response.CouponInfo {
	discVal, _ := strconv.ParseFloat(c.DiscountValue, 64)
	minOrder, _ := strconv.ParseFloat(c.MinOrderAmount, 64)
	return response.CouponInfo{
		ID:             c.ID,
		Code:           c.Code,
		DiscountType:   c.DiscountType,
		DiscountValue:  discVal,
		MinOrderAmount: minOrder,
		MaxUses:        c.MaxUses,
		UsesRemaining:  c.UsesRemaining,
		EventID:        c.EventID,
		ValidFrom:      c.ValidFrom,
		ValidUntil:     c.ValidUntil,
	}
}

func ToPricingRuleInfo(r *models.PricingRule) response.PricingRuleInfo {
	mult, _ := strconv.ParseFloat(r.Multiplier, 64)
	return response.PricingRuleInfo{
		ID:         r.ID,
		RuleType:   r.RuleType,
		LowerBound: r.LowerBound,
		UpperBound: r.UpperBound,
		Multiplier: mult,
		Active:     r.Active,
	}
}
