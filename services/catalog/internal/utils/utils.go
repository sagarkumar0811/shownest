package utils

import (
	"github.com/gin-gonic/gin"
	pkgutils "github.com/shownest/pkg/utils"
)

func NewUUID() (string, error) {
	return pkgutils.NewUUID()
}

func GetS3Key(service string, parts ...string) string {
	return pkgutils.GetS3Key(service, parts...)
}

func JoinColumns(columns []string) string {
	return pkgutils.JoinColumns(columns)
}

func MustUserID(c *gin.Context) string {
	return pkgutils.MustUserID(c)
}

func WriteError(c *gin.Context, err error) {
	pkgutils.WriteError(c, err)
}
