package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/user-service/internal/dto/request"
	"github.com/shownest/user-service/internal/utils"
)

func (h *Handler) ListSessions(c *gin.Context) {
	sessions, err := h.usecase.ListSessions(c.Request.Context(), utils.MustUserID(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessions": sessions})
}

func (h *Handler) RevokeSession(c *gin.Context) {
	var req request.RevokeSessionRequest
	if err := c.ShouldBindUri(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, "session id is required"))
		return
	}
	if err := h.usecase.RevokeSession(c.Request.Context(), utils.MustUserID(c), req.SessionID); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session revoked"})
}

func (h *Handler) RevokeAllSessions(c *gin.Context) {
	if err := h.usecase.RevokeAllSessions(c.Request.Context(), utils.MustUserID(c), utils.MustSessionID(c)); err != nil {
		utils.WriteError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all other sessions revoked"})
}
