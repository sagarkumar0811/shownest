package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/logger"
	"github.com/shownest/pkg/utils"
	"github.com/shownest/pricing-service/internal/dto/request"
	"go.uber.org/zap"
)

func (h *Handler) CalculatePrice(c *gin.Context) {
	var req request.CalculatePriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	result, err := h.usecase.CalculatePrice(c.Request.Context(), utils.MustUserID(c), req.ShowtimeID, req.SeatIDs, req.CouponCode)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("calculate price failed", zap.String("showtimeId", req.ShowtimeID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ValidateCoupon(c *gin.Context) {
	var req request.ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	result, err := h.usecase.ValidateCoupon(c.Request.Context(), utils.MustUserID(c), req.Code, req.ShowtimeID, req.OrderTotal)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("validate coupon failed", zap.String("code", req.Code), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
