package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/merchant-service/internal/dto/request"
	"github.com/shownest/merchant-service/internal/utils"
	apperrors "github.com/shownest/pkg/errors"
)

func (h *Handler) RequestDocumentUploadURL(c *gin.Context) {
	var req request.DocumentUploadURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	resp, err := h.usecase.RequestDocumentUploadURL(c.Request.Context(), utils.MustUserID(c), req.DocumentType)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) ConfirmDocument(c *gin.Context) {
	var req request.ConfirmDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	if err := req.Validate(); err != nil {
		utils.WriteError(c, apperrors.New(apperrors.CodeInvalidArgument, err.Error()))
		return
	}
	info, err := h.usecase.ConfirmDocument(c.Request.Context(), utils.MustUserID(c), req)
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, info)
}

func (h *Handler) ListDocuments(c *gin.Context) {
	docs, err := h.usecase.ListDocuments(c.Request.Context(), utils.MustUserID(c))
	if err != nil {
		utils.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"documents": docs})
}
