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

func (h *Handler) RequestMediaUploadURL(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.MediaUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	resp, err := h.usecase.RequestMediaUploadURL(c.Request.Context(), utils.MustUserID(c), uri.EventID, req.MediaType)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("request media upload url failed", zap.String("eventId", uri.EventID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ConfirmMedia(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.ConfirmMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.ConfirmMedia(c.Request.Context(), utils.MustUserID(c), uri.EventID, req)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("confirm media failed", zap.String("eventId", uri.EventID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ListMedia(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.ListMedia(c.Request.Context(), uri.EventID)
	if err != nil {
		logger.WithContext(c.Request.Context()).Error("list media failed", zap.String("eventId", uri.EventID), zap.Error(err))
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}
