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
	pkgutils "github.com/shownest/pkg/utils"
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

func MustSessionID(c *gin.Context) string {
	return c.MustGet("sessionId").(string)
}

func UserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func NewUUID() (string, error)             { return pkgutils.NewUUID() }
func JoinColumns(columns []string) string  { return pkgutils.JoinColumns(columns) }
func MustUserID(c *gin.Context) string     { return pkgutils.MustUserID(c) }
func ClientIP(c *gin.Context) string       { return pkgutils.ClientIP(c) }
func WriteError(c *gin.Context, err error) { pkgutils.WriteError(c, err) }
