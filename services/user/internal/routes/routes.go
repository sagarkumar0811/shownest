package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/pkg/jwt"
	"github.com/shownest/pkg/middleware"
	"github.com/shownest/user-service/internal/handlers"
)

type Config struct {
	Handler    *handlers.Handler
	JWTService *jwt.Service
	Port       string
}

func InitRoutes(config Config) error {
	r := gin.New()
	r.Use(gin.Recovery())

	base := r.Group("/api/user")

	// Health check endpoint
	base.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})

	// API v1 routes
	v1 := base.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/otp/send", config.Handler.SendOTP)
			auth.POST("/otp/verify", config.Handler.VerifyOTP)
			auth.POST("/token/refresh", config.Handler.RefreshToken)

			protected := auth.Group("", middleware.JWTAuth(config.JWTService))
			{
				protected.GET("/sessions", config.Handler.ListSessions)
				protected.DELETE("/sessions/:id", config.Handler.RevokeSession)
				protected.DELETE("/sessions", config.Handler.RevokeAllSessions)
			}
		}

		me := v1.Group("/me", middleware.JWTAuth(config.JWTService))
		{
			me.GET("", config.Handler.GetProfile)
			me.PATCH("", config.Handler.UpdateProfile)
		}
	}

	addr := fmt.Sprintf(":%s", config.Port)
	return r.Run(addr)
}
