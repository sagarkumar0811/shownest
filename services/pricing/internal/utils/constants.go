package utils

import "time"

const (
	PricingRulesCacheKey = "pricing:rules"
	PricingCacheTTL      = 5 * time.Minute

	RuleTypeTimeOfDay = "hourly"
	RuleTypeDayOfWeek = "weekly"
	RuleTypeSurge     = "surge"

	DiscountTypeFlat    = "flat"
	DiscountTypePercent = "percent"
)
