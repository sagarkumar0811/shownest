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

func (h *Handler) CreateEvent(c *gin.Context) {
	var req request.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.CreateEvent(c.Request.Context(), utils.MustUserID(c), req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("create event failed", zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) GetEvent(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.GetEvent(c.Request.Context(), uri.EventID)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("get event failed", zap.String("eventId", uri.EventID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) ListEvents(c *gin.Context) {
	var q request.ListEventsRequest
	if err := c.ShouldBindQuery(&q); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.ListEvents(c.Request.Context(), q)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("list events failed", zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}

func (h *Handler) UpdateEvent(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.UpdateEvent(c.Request.Context(), utils.MustUserID(c), uri.EventID, req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("update event failed", zap.String("eventId", uri.EventID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}
