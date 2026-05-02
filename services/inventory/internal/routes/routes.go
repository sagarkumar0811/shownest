package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/inventory-service/internal/handlers"
	"github.com/shownest/pkg/jwt"
	"github.com/shownest/pkg/middleware"
)

type Config struct {
	Handler    *handlers.Handler
	JWTService *jwt.Service
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
		v1.GET("/showtimes/:showtimeId/seats", config.Handler.ListShowtimeSeats)
		v1.GET("/halls/:hallId/seats", config.Handler.ListSeats)
		v1.GET("/halls/:hallId/seat-categories", config.Handler.ListSeatCategories)

		// Authenticated user routes
		user := v1.Group("", auth)
		user.POST("/showtimes/:showtimeId/seats/lock", config.Handler.LockSeats)

		// Merchant-only routes
		merchant := v1.Group("", auth, merchantOnly)
		{
			merchant.POST("/halls/:hallId/seat-categories", config.Handler.CreateSeatCategory)
			merchant.POST("/halls/:hallId/seats", config.Handler.BulkCreateSeats)
			merchant.POST("/halls/:hallId/showtimes/publish", config.Handler.PublishShowtimeSeats)
		}
	}

	// Internal S2S routes
	internal := base.Group("/internal")
	{
		internal.POST("/showtimes/:showtimeId/seats/confirm", config.Handler.ConfirmSeats)
		internal.POST("/showtimes/:showtimeId/seats/release", config.Handler.ReleaseSeats)
		internal.POST("/showtimes/:showtimeId/seats/prices", config.Handler.GetSeatPrices)
	}

	addr := fmt.Sprintf(":%s", config.Port)
	return r.Run(addr)
}
