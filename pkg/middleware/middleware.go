package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	apperrors "github.com/shownest/pkg/errors"
	"github.com/shownest/pkg/jwt"
)

func JWTAuth(jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			abortWithError(c, apperrors.New(apperrors.CodeUnauthenticated, "authorization header required"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			abortWithError(c, apperrors.New(apperrors.CodeUnauthenticated, "authorization header format invalid"))
			return
		}

		claims, err := jwtService.ValidateAccessToken(parts[1])
		if err != nil {
			switch err {
			case jwt.ErrTokenExpired:
				abortWithError(c, apperrors.New(apperrors.CodeTokenExpired, "access token expired"))
			default:
				abortWithError(c, apperrors.New(apperrors.CodeTokenInvalid, "invalid access token"))
			}
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("sessionId", claims.SessionID)
		c.Next()
	}
}

func abortWithError(c *gin.Context, appErr *apperrors.AppError) {
	c.AbortWithStatusJSON(apperrors.HTTPStatus(appErr.Code), gin.H{
		"error": gin.H{
			"code":    appErr.Code.String(),
			"message": appErr.Message,
		},
	})
}
