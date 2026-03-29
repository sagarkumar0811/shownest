package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/merchant-service/internal/dto/request"
	"github.com/shownest/merchant-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
)

func (h *Handler) CreateHall(c *gin.Context) {
	var uriReq request.VenueIDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, "venue id is required"))
		return
	}
	var req request.CreateHallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.CreateHall(c.Request.Context(), uriReq.VenueID, utils.MustUserID(c), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ListHalls(c *gin.Context) {
	var uriReq request.VenueIDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, "venue id is required"))
		return
	}
	halls, err := h.usecase.ListHalls(c.Request.Context(), uriReq.VenueID, utils.MustUserID(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"halls": halls})
}
