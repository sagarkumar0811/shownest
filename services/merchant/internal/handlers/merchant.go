package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/merchant-service/internal/dto/request"
	"github.com/shownest/merchant-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
)

func (h *Handler) CreateMerchant(c *gin.Context) {
	var req request.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.CreateMerchant(c.Request.Context(), utils.MustUserID(c), req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("create merchant failed", zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) GetMyMerchant(c *gin.Context) {
	info, err := h.usecase.GetMyMerchant(c.Request.Context(), utils.MustUserID(c))
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("get my merchant failed", zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) GetMerchantByUserID(c *gin.Context) {
	userID := c.Param("userId")
	info, err := h.usecase.GetMyMerchant(c.Request.Context(), userID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": info.ID})
}

func (h *Handler) SubmitForReview(c *gin.Context) {
	if err := h.usecase.SubmitForReview(c.Request.Context(), utils.MustUserID(c)); err != nil {
		logger.WithContext(c.Request.Context()).Error("submit for review failed", zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "submitted for review"})
}
