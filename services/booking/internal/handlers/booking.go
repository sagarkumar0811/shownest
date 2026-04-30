package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/booking-service/internal/dto/request"
	apperrors "github.com/shownest/pkg/errors"
	pkgutils "github.com/shownest/pkg/utils"
)

func (h *Handler) CreateBooking(c *gin.Context) {
	var req request.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	userID := pkgutils.MustUserID(c)
	info, err := h.usecase.CreateBooking(c.Request.Context(), userID, req)
	if err != nil {
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ConfirmBooking(c *gin.Context) {
	var uri request.BookingIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	userID := pkgutils.MustUserID(c)

	info, err := h.usecase.ConfirmBooking(c.Request.Context(), uri.ID, userID)
	if err != nil {
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	var uri request.BookingIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	userID := pkgutils.MustUserID(c)

	if err := h.usecase.CancelBooking(c.Request.Context(), uri.ID, userID); err != nil {
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}

func (h *Handler) GetBooking(c *gin.Context) {
	var uri request.BookingIDRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		pkgutils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	userID := pkgutils.MustUserID(c)
	info, err := h.usecase.GetBooking(c.Request.Context(), uri.ID, userID)
	if err != nil {
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, info)
}

func (h *Handler) ListUserBookings(c *gin.Context) {
	userID := pkgutils.MustUserID(c)
	infos, err := h.usecase.ListUserBookings(c.Request.Context(), userID)
	if err != nil {
		pkgutils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, infos)
}
