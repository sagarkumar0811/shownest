package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
)

func FilterMobileNumber(phone string) (string, error) {

	// strip whitespace
	rr := make([]rune, 0, len(phone))
	for _, r := range phone {
		if !unicode.IsSpace(r) {
			rr = append(rr, r)
		}
	}
	trimmed := string(rr)

	// resolve scientific notation (e.g. 9.1234e9)
	if strings.ContainsAny(trimmed, "eE") {
		f, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return phone, err
		}
		trimmed = fmt.Sprintf("%.0f", f)
	}

	flt, _, err := big.ParseFloat(trimmed, 10, 0, big.ToNearestEven)
	if err != nil {
		return phone, err
	}
	k := new(big.Int)
	k, _ = flt.Int(k)

	re := regexp.MustCompile(`[0-9]+`)
	s := re.FindAllString(k.String(), -1)
	if len(s) == 0 {
		return "", errors.New("invalid phone format")
	}

	str := s[0]
	if len(str) < 10 || len(str) > 12 {
		return str, errors.New("invalid phone length")
	}

	last10 := str[len(str)-10:]
	if match, _ := regexp.MatchString(`^[6789][0-9]{9}$`, last10); !match {
		return str, errors.New("invalid phone number")
	}

	return last10, nil
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

func NewUUID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
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
