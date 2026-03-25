package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
)

func IsValidPhone(phone string) bool {
	e164Regex := regexp.MustCompile(`^\+[1-9]\d{7,14}$`)
	return e164Regex.MatchString(phone)
}

func GenerateOTP() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	n := int(b[0])<<16 | int(b[1])<<8 | int(b[2])
	return fmt.Sprintf("%06d", n%1_000_000), nil
}

func HashSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

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

func JoinColumns(columns []string) string {
	var result strings.Builder
	for i, c := range columns {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(c)
	}
	return result.String()
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
