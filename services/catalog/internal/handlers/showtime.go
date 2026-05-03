package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/catalog-service/internal/dto/request"
	"github.com/shownest/catalog-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/logger"
	"go.uber.org/zap"
)

func (h *Handler) CreateShowtime(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.CreateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.CreateShowtime(c.Request.Context(), utils.MustUserID(c), uri.EventID, req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("create showtime failed", zap.String("eventId", uri.EventID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) GetShowtime(c *gin.Context) {
	var uri request.ShowtimeIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.GetShowtime(c.Request.Context(), uri.ShowtimeID)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("get showtime failed", zap.String("showtimeId", uri.ShowtimeID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) GetShowtimeBasePrice(c *gin.Context) {
	var uri request.ShowtimeIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	price, err := h.usecase.GetShowtimeBasePrice(c.Request.Context(), uri.ShowtimeID)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("get showtime base price failed", zap.String("showtimeId", uri.ShowtimeID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"basePrice": price})
}

func (h *Handler) ListShowtimes(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.ListShowtimes(c.Request.Context(), uri.EventID)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("list showtimes failed", zap.String("eventId", uri.EventID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}

func (h *Handler) UpdateShowtime(c *gin.Context) {
	var uri request.ShowtimeIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.UpdateShowtimeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if req.Status == nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, "status is required"))
		return
	}
	info, err := h.usecase.UpdateShowtime(c.Request.Context(), utils.MustUserID(c), uri.ShowtimeID, req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("update showtime failed", zap.String("showtimeId", uri.ShowtimeID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}
