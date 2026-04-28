package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/catalog-service/internal/dto/request"
	"github.com/shownest/catalog-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
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
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) ListShowtimes(c *gin.Context) {
	var uri request.EventIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.ListShowtimes(c.Request.Context(), uri.EventID)
	if err != nil {
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
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}
