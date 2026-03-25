package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/user-service/internal/dto/request"
	"github.com/shownest/user-service/internal/utils"
)

func (h *Handler) GetProfile(c *gin.Context) {
	profile, err := h.usecase.GetProfile(c.Request.Context(), utils.MustUserID(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}

	profile, err := h.usecase.UpdateProfile(c.Request.Context(), utils.MustUserID(c), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, profile)
}
