package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
)

func NewUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func NewHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func GetS3Key(service string, parts ...string) string {
	return service + "/" + strings.Join(parts, "/")
}

func JoinColumns(columns []string) string {
	var sb strings.Builder
	for i, c := range columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(c)
	}
	return sb.String()
}

func MustUserID(c *gin.Context) string {
	return c.MustGet("userId").(string)
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
