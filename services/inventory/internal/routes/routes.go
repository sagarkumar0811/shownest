package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shownest/inventory-service/internal/handlers"
	"github.com/shownest/pkg/jwt"
	"github.com/shownest/pkg/middleware"
)

type Config struct {
	Handler    *handlers.Handler
	JWTService *jwt.Service
	Cache      *redis.Client
	Port       string
}

func InitRoutes(config Config) error {
	r := gin.New()
	r.Use(gin.Recovery())

	base := r.Group("/api/inventory")

	base.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})

	v1 := base.Group("/v1")
	{
		auth := middleware.JWTAuth(config.JWTService)
		merchantOnly := middleware.RequireMerchant()

		// Public routes
		v1.GET("/showtimes/:id/seats", config.Handler.ListShowtimeSeats)
		v1.GET("/halls/:id/seats", config.Handler.ListSeats)
		v1.GET("/halls/:id/seat-categories", config.Handler.ListSeatCategories)

		// Authenticated user routes
		user := v1.Group("", auth)
		user.POST("/showtimes/:id/seats/lock", config.Handler.LockSeats)

		// Merchant-only routes
		merchant := v1.Group("", auth, merchantOnly)
		{
			merchant.POST("/halls/:id/seat-categories", config.Handler.CreateSeatCategory)
			merchant.POST("/halls/:id/seats", config.Handler.BulkCreateSeats)
			merchant.POST("/halls/:id/showtimes/publish", config.Handler.PublishShowtimeSeats)
		}
	}

	// Internal S2S routes
	internal := base.Group("/internal")
	{
		internal.POST("/showtimes/:id/seats/confirm", config.Handler.ConfirmSeats)
		internal.POST("/showtimes/:id/seats/release", config.Handler.ReleaseSeats)
	}

	addr := fmt.Sprintf(":%s", config.Port)
	return r.Run(addr)
}
