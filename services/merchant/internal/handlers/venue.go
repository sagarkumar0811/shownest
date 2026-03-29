package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/merchant-service/internal/dto/request"
	"github.com/shownest/merchant-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
)

func (h *Handler) CreateVenue(c *gin.Context) {
	var req request.CreateVenueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.CreateVenue(c.Request.Context(), utils.MustUserID(c), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ListMyVenues(c *gin.Context) {
	venues, err := h.usecase.ListMyVenues(c.Request.Context(), utils.MustUserID(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"venues": venues})
}

func (h *Handler) GetVenue(c *gin.Context) {
	var req request.VenueIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, "venue id is required"))
		return
	}
	info, err := h.usecase.GetVenue(c.Request.Context(), req.VenueID, utils.MustUserID(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) GetNearbyVenues(c *gin.Context) {
	var req request.NearbyVenuesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	venues, err := h.usecase.GetNearbyVenues(c.Request.Context(), req.Latitude, req.Longitude, req.Radius)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"venues": venues})
}

func (h *Handler) GetVenuesByCity(c *gin.Context) {
	var req request.CityRequest
	if err := c.ShouldBindUri(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, "city is required"))
		return
	}
	venues, err := h.usecase.GetVenuesByCity(c.Request.Context(), req.City)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"venues": venues})
}
