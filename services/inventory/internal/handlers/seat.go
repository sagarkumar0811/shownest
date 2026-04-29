package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/inventory-service/internal/dto/request"
	"github.com/shownest/inventory-service/internal/utils"
)

func (h *Handler) CreateSeatCategory(c *gin.Context) {
	var uri request.HallIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.CreateSeatCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.CreateSeatCategory(c.Request.Context(), req, uri.HallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ListSeatCategories(c *gin.Context) {
	var uri request.HallIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.ListSeatCategories(c.Request.Context(), uri.HallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}

func (h *Handler) BulkCreateSeats(c *gin.Context) {
	var uri request.HallIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.BulkCreateSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.BulkCreateSeats(c.Request.Context(), uri.HallID, req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, infos)
}

func (h *Handler) ListSeats(c *gin.Context) {
	var uri request.HallIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.ListSeats(c.Request.Context(), uri.HallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}

func (h *Handler) PublishShowtimeSeats(c *gin.Context) {
	var req request.PublishShowtimeSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var uri request.HallIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	n, err := h.usecase.PublishShowtimeSeats(c.Request.Context(), req.ShowtimeID, uri.HallID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"seatsCreated": n})
}

func (h *Handler) ListShowtimeSeats(c *gin.Context) {
	var uri request.ShowtimeIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	infos, err := h.usecase.ListShowtimeSeats(c.Request.Context(), uri.ShowtimeID)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}

func (h *Handler) LockSeats(c *gin.Context) {
	var uri request.ShowtimeIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.LockSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	result, err := h.usecase.LockSeats(c.Request.Context(), uri.ShowtimeID, utils.MustUserID(c), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ReleaseSeats(c *gin.Context) {
	var uri request.ShowtimeIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	var req request.ReleaseSeatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := h.usecase.ReleaseSeats(c.Request.Context(), uri.ShowtimeID, utils.MustUserID(c), req); err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"released": true})
}
