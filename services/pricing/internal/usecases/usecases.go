package usecases

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/logger"
	"github.com/shownest/pricing-service/internal/client"
	"github.com/shownest/pricing-service/internal/dto/request"
	"github.com/shownest/pricing-service/internal/dto/response"
	"github.com/shownest/pricing-service/internal/mapper"
	"github.com/shownest/pricing-service/internal/models"
	"github.com/shownest/pricing-service/internal/repository"
	"github.com/shownest/pricing-service/internal/utils"
	"go.uber.org/zap"
)

type UseCase struct {
	repo      *repository.Repository
	cache     *redis.Client
	catalog   *client.CatalogClient
	inventory *client.InventoryClient
	merchant  *client.MerchantClient
}

func New(repo *repository.Repository, cache *redis.Client, catalog *client.CatalogClient, inventory *client.InventoryClient, merchant *client.MerchantClient) *UseCase {
	return &UseCase{repo: repo, cache: cache, catalog: catalog, inventory: inventory, merchant: merchant}
}

func (uc *UseCase) CalculatePrice(ctx context.Context, userID, showtimeID string, seatIDs []string, couponCode *string) (*response.PriceCalculationResponse, error) {
	showtime, err := uc.catalog.GetShowtime(ctx, showtimeID)
	if err != nil {
		logger.WithContext(ctx).Error("fetch showtime for pricing failed", zap.String("showtimeId", showtimeID), zap.Error(err))
		return nil, err
	}

	event, err := uc.catalog.GetEvent(ctx, showtime.EventID)
	if err != nil {
		logger.WithContext(ctx).Error("fetch event for pricing failed", zap.String("eventId", showtime.EventID), zap.Error(err))
		return nil, err
	}

	seatPrices, err := uc.inventory.GetSeatPrices(ctx, showtimeID, seatIDs)
	if err != nil {
		logger.WithContext(ctx).Error("fetch seat prices failed", zap.String("showtimeId", showtimeID), zap.Error(err))
		return nil, err
	}

	occupancy, err := uc.inventory.GetOccupancy(ctx, showtimeID)
	if err != nil {
		logger.WithContext(ctx).Error("fetch occupancy failed", zap.String("showtimeId", showtimeID), zap.Error(err))
		return nil, err
	}

	rules, err := uc.loadPricingRulesByMerchant(ctx, event.MerchantID)
	if err != nil {
		logger.WithContext(ctx).Error("load pricing rules failed", zap.Error(err))
		return nil, err
	}

	basePrice, _ := strconv.ParseFloat(showtime.BasePrice, 64)
	seatBreakdowns, subtotal := mapper.ToSeatPriceBreakdownList(basePrice, seatPrices)
	timeMultiplier, dayMultiplier, surgeMultiplier := applyDynamicMultipliers(rules, showtime, occupancy.OccupancyPercent)

	dynamicSubtotal := subtotal * timeMultiplier * dayMultiplier * surgeMultiplier

	var couponDiscount float64
	if couponCode != nil && *couponCode != "" {
		couponDiscount = uc.computeCouponDiscount(ctx, userID, *couponCode, dynamicSubtotal, showtime.EventID)
	}

	total := dynamicSubtotal - couponDiscount
	if total < 0 {
		total = 0
	}

	return &response.PriceCalculationResponse{
		ShowtimeID:      showtimeID,
		Seats:           seatBreakdowns,
		Subtotal:        subtotal,
		TimeMultiplier:  timeMultiplier,
		DayMultiplier:   dayMultiplier,
		SurgeMultiplier: surgeMultiplier,
		DynamicSubtotal: dynamicSubtotal,
		CouponDiscount:  couponDiscount,
		Total:           total,
	}, nil
}

func (uc *UseCase) ValidateCoupon(ctx context.Context, userID, code, showtimeID string, orderTotal float64) (*response.ValidateCouponResponse, error) {
	showtime, err := uc.catalog.GetShowtime(ctx, showtimeID)
	if err != nil {
		logger.WithContext(ctx).Error("fetch showtime for coupon validation failed", zap.String("showtimeId", showtimeID), zap.Error(err))
		return nil, err
	}

	coupon, err := uc.repo.GetCouponByCode(ctx, code)
	if err != nil {
		if apperrors.HasCode(err, apperrors.CodeDBNotFound) {
			return &response.ValidateCouponResponse{Valid: false, Message: "coupon not found"}, nil
		}
		logger.WithContext(ctx).Error("get coupon by code failed", zap.String("code", code), zap.Error(err))
		return nil, err
	}

	now := time.Now()
	if now.Before(coupon.ValidFrom) || now.After(coupon.ValidUntil) {
		return &response.ValidateCouponResponse{Valid: false, Message: "coupon is not valid at this time"}, nil
	}

	if coupon.UsesRemaining != nil && *coupon.UsesRemaining <= 0 {
		return &response.ValidateCouponResponse{Valid: false, Message: "coupon has no uses remaining"}, nil
	}

	minOrder, _ := strconv.ParseFloat(coupon.MinOrderAmount, 64)
	if orderTotal < minOrder {
		return &response.ValidateCouponResponse{Valid: false, Message: "order total is below minimum for this coupon"}, nil
	}

	if coupon.EventID != nil && *coupon.EventID != showtime.EventID {
		return &response.ValidateCouponResponse{Valid: false, Message: "coupon is not valid for this event"}, nil
	}

	_, err = uc.repo.GetRedemptionByUser(ctx, coupon.ID, userID)
	if err == nil {
		return &response.ValidateCouponResponse{Valid: false, Message: "coupon already used by this account"}, nil
	}
	if !apperrors.HasCode(err, apperrors.CodeDBNotFound) {
		logger.WithContext(ctx).Error("check coupon redemption failed", zap.String("couponId", coupon.ID), zap.Error(err))
		return nil, err
	}

	discount := uc.calculateDiscount(coupon, orderTotal)
	info := mapper.ToCouponInfo(coupon)

	return &response.ValidateCouponResponse{
		Valid:          true,
		Coupon:         &info,
		DiscountAmount: discount,
	}, nil
}

func (uc *UseCase) CreateCoupon(ctx context.Context, userID string, req request.CreateCouponRequest) (*response.CouponInfo, error) {
	merchantID, err := uc.merchant.GetMerchantIDByUserID(ctx, userID)
	if err != nil {
		logger.WithContext(ctx).Error("resolve merchant id failed", zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	coupon, err := uc.repo.CreateCoupon(ctx,
		merchantID, req.Code, req.DiscountType, req.DiscountValue,
		req.MinOrderAmount, req.MaxUses, req.EventID, req.ValidFrom, req.ValidUntil,
	)
	if err != nil {
		logger.WithContext(ctx).Error("create coupon failed", zap.String("merchantId", merchantID), zap.Error(err))
		return nil, err
	}
	info := mapper.ToCouponInfo(coupon)
	return &info, nil
}

func (uc *UseCase) ListMyCoupons(ctx context.Context, userID string) ([]response.CouponInfo, error) {
	merchantID, err := uc.merchant.GetMerchantIDByUserID(ctx, userID)
	if err != nil {
		logger.WithContext(ctx).Error("resolve merchant id failed", zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	coupons, err := uc.repo.ListCouponsByMerchant(ctx, merchantID)
	if err != nil {
		logger.WithContext(ctx).Error("list coupons failed", zap.String("merchantId", merchantID), zap.Error(err))
		return nil, err
	}
	infos := make([]response.CouponInfo, len(coupons))
	for i, c := range coupons {
		infos[i] = mapper.ToCouponInfo(&c)
	}
	return infos, nil
}

func (uc *UseCase) RecordRedemption(ctx context.Context, couponCode, userID, bookingID string) error {
	coupon, err := uc.repo.GetCouponByCode(ctx, couponCode)
	if err != nil {
		logger.WithContext(ctx).Error("fetch coupon by code failed", zap.String("couponCode", couponCode), zap.Error(err))
		return err
	}

	if _, err := uc.repo.DecrementCouponUses(ctx, coupon.ID); err != nil && !apperrors.HasCode(err, apperrors.CodeFailedPrecondition) {
		logger.WithContext(ctx).Error("decrement coupon uses failed", zap.String("couponCode", couponCode), zap.Error(err))
		return err
	}

	if _, err := uc.repo.CreateRedemption(ctx, coupon.ID, userID, bookingID); err != nil {
		logger.WithContext(ctx).Error("create redemption failed", zap.String("couponCode", couponCode), zap.Error(err))
		return err
	}
	return nil
}

func (uc *UseCase) CreatePricingRule(ctx context.Context, userID string, req request.CreatePricingRuleRequest) (*response.PricingRuleInfo, error) {
	merchantID, err := uc.merchant.GetMerchantIDByUserID(ctx, userID)
	if err != nil {
		logger.WithContext(ctx).Error("resolve merchant id failed", zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	rule, err := uc.repo.CreatePricingRule(ctx, merchantID, req.RuleType, req.LowerBound, req.UpperBound, req.Multiplier)
	if err != nil {
		logger.WithContext(ctx).Error("create pricing rule failed", zap.String("ruleType", req.RuleType), zap.Error(err))
		return nil, err
	}

	uc.cache.Del(ctx, utils.PricingRulesCacheKey(merchantID))

	info := mapper.ToPricingRuleInfo(rule)
	return &info, nil
}

func (uc *UseCase) ListPricingRules(ctx context.Context, userID string) ([]response.PricingRuleInfo, error) {
	merchantID, err := uc.merchant.GetMerchantIDByUserID(ctx, userID)
	if err != nil {
		logger.WithContext(ctx).Error("resolve merchant id failed", zap.String("userId", userID), zap.Error(err))
		return nil, err
	}

	rules, err := uc.repo.ListActivePricingRulesByMerchant(ctx, merchantID)
	if err != nil {
		logger.WithContext(ctx).Error("list pricing rules failed", zap.String("merchantId", merchantID), zap.Error(err))
		return nil, err
	}
	infos := make([]response.PricingRuleInfo, len(rules))
	for i, r := range rules {
		infos[i] = mapper.ToPricingRuleInfo(&r)
	}
	return infos, nil
}

func (uc *UseCase) loadPricingRulesByMerchant(ctx context.Context, merchantID string) ([]models.PricingRule, error) {
	cacheKey := utils.PricingRulesCacheKey(merchantID)
	if cached, err := uc.cache.Get(ctx, cacheKey).Bytes(); err == nil {
		var rules []models.PricingRule
		if json.Unmarshal(cached, &rules) == nil {
			return rules, nil
		}
	}

	rules, err := uc.repo.ListActivePricingRulesByMerchant(ctx, merchantID)
	if err != nil {
		return nil, err
	}

	if b, err := json.Marshal(rules); err == nil {
		uc.cache.Set(ctx, cacheKey, b, utils.PricingCacheTTL)
	}
	return rules, nil
}

func applyDynamicMultipliers(rules []models.PricingRule, showtime *client.ShowtimeInfo, occupancyPercent float64) (timeMultiplier, dayMultiplier, surgeMultiplier float64) {
	timeMultiplier, dayMultiplier, surgeMultiplier = 1.0, 1.0, 1.0
	hour := showtime.StartTime.Hour()
	dayOfWeek := int(showtime.StartTime.Weekday())

	for _, rule := range rules {
		multiplier, _ := strconv.ParseFloat(rule.Multiplier, 64)
		switch rule.RuleType {
		case utils.RuleTypeTimeOfDay:
			if rule.LowerBound != nil && rule.UpperBound != nil &&
				hour >= *rule.LowerBound && hour < *rule.UpperBound {
				timeMultiplier = multiplier
			}
		case utils.RuleTypeDayOfWeek:
			if rule.LowerBound != nil && dayOfWeek == *rule.LowerBound {
				dayMultiplier = multiplier
			}
		case utils.RuleTypeSurge:
			if rule.LowerBound != nil && occupancyPercent >= float64(*rule.LowerBound) {
				surgeMultiplier = multiplier
			}
		}
	}
	return
}

func (uc *UseCase) computeCouponDiscount(ctx context.Context, userID, code string, orderTotal float64, eventID string) float64 {
	coupon, err := uc.repo.GetCouponByCode(ctx, code)
	if err != nil {
		return 0
	}
	now := time.Now()
	if now.Before(coupon.ValidFrom) || now.After(coupon.ValidUntil) {
		return 0
	}
	if coupon.UsesRemaining != nil && *coupon.UsesRemaining <= 0 {
		return 0
	}
	if coupon.EventID != nil && *coupon.EventID != eventID {
		return 0
	}
	_, err = uc.repo.GetRedemptionByUser(ctx, coupon.ID, userID)
	if err == nil {
		return 0
	}
	minOrder, _ := strconv.ParseFloat(coupon.MinOrderAmount, 64)
	if orderTotal < minOrder {
		return 0
	}
	return uc.calculateDiscount(coupon, orderTotal)
}

func (uc *UseCase) calculateDiscount(coupon *models.Coupon, orderTotal float64) float64 {
	discVal, _ := strconv.ParseFloat(coupon.DiscountValue, 64)
	if coupon.DiscountType == utils.DiscountTypeFlat {
		if discVal > orderTotal {
			return orderTotal
		}
		return discVal
	}
	return orderTotal * discVal / 100
}
