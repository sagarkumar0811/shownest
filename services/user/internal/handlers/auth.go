package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/user-service/internal/dto/request"
	"github.com/shownest/user-service/internal/utils"
)

func (h *Handler) SendOTP(c *gin.Context) {
	var req request.SendOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	if err := h.usecase.SendOTP(c.Request.Context(), req.Phone); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (h *Handler) VerifyOTP(c *gin.Context) {
	var req request.VerifyOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	resp, err := h.usecase.VerifyOTP(c.Request.Context(), req, utils.UserAgent(c), utils.ClientIP(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req request.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	resp, err := h.usecase.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
