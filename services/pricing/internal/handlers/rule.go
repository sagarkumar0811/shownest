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

func (h *Handler) CreatePricingRule(c *gin.Context) {
	var req request.CreatePricingRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	info, err := h.usecase.CreatePricingRule(c.Request.Context(), pkgutils.MustUserID(c), req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("create pricing rule failed", zap.String("ruleType", req.RuleType), zap.Error(err))
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ListPricingRules(c *gin.Context) {
	infos, err := h.usecase.ListPricingRules(c.Request.Context(), pkgutils.MustUserID(c))
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("list pricing rules failed", zap.Error(err))
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"rules": infos})
}
