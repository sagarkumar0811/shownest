package utils

import "time"

const (
	pricingRulesPrefix = "pricing:rules:"
	PricingCacheTTL    = 5 * time.Minute
)

func PricingRulesCacheKey(merchantID string) string {
	return pricingRulesPrefix + merchantID
}
