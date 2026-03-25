package utils

import (
	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
)

func MustUserID(c *gin.Context) string {
	return c.MustGet("userId").(string)
}

func MustSessionID(c *gin.Context) string {
	return c.MustGet("sessionId").(string)
}

func UserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func ClientIP(c *gin.Context) string {
	return c.ClientIP()
}

func WriteError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if !apperrors.As(err, &appErr) {
		appErr = apperrors.New(apperrors.CodeInternal, "internal server error")
	}
	c.AbortWithStatusJSON(apperrors.HTTPStatus(appErr.Code), gin.H{
		"error": gin.H{
			"code":    appErr.Code.String(),
			"message": appErr.Message,
		},
	})
}
