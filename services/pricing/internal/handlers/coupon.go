package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/logger"
	pkgutils "github.com/shownest/pkg/utils"
	"github.com/shownest/pricing-service/internal/dto/request"
	"go.uber.org/zap"
)

func (h *Handler) CreateCoupon(c *gin.Context) {
	var req request.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	userID := pkgutils.MustUserID(c)

	info, err := h.usecase.CreateCoupon(c.Request.Context(), userID, req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("create coupon failed", zap.String("code", req.Code), zap.Error(err))
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ListMyCoupons(c *gin.Context) {
	userID := pkgutils.MustUserID(c)

	infos, err := h.usecase.ListMyCoupons(c.Request.Context(), userID)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("list coupons failed", zap.String("userId", userID), zap.Error(err))
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"coupons": infos})
}

func (h *Handler) RecordRedemption(c *gin.Context) {
	couponCode := c.Param("couponCode")

	var req request.RecordRedemptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	if err := h.usecase.RecordRedemption(c.Request.Context(), couponCode, req.UserID, req.BookingID); err != nil {
		logger.WithContext(c.Request.Context()).Error("record redemption failed", zap.String("couponCode", couponCode), zap.Error(err))
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"recorded": true})
}
