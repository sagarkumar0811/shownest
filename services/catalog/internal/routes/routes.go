package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/shownest/catalog-service/internal/handlers"
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

	base := r.Group("/api/catalog")

	base.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})

	v1 := base.Group("/v1")
	{
		auth := middleware.JWTAuth(config.JWTService)
		merchantOnly := middleware.RequireMerchant()

		// Public routes
		v1.GET("/events", config.Handler.ListEvents)
		v1.GET("/events/:eventId", config.Handler.GetEvent)
		v1.GET("/events/:eventId/showtimes", config.Handler.ListShowtimes)
		v1.GET("/events/:eventId/media", config.Handler.ListMedia)
		v1.GET("/showtimes/:showtimeId", config.Handler.GetShowtime)

		// Merchant-only routes
		merchant := v1.Group("", auth, merchantOnly)
		{
			merchant.POST("/events", config.Handler.CreateEvent)
			merchant.PATCH("/events/:eventId", config.Handler.UpdateEvent)
			merchant.POST("/events/:eventId/showtimes", config.Handler.CreateShowtime)
			merchant.PATCH("/showtimes/:showtimeId", config.Handler.UpdateShowtime)
			merchant.POST("/events/:eventId/media/upload-url", config.Handler.RequestMediaUploadURL)
			merchant.POST("/events/:eventId/media/confirm", config.Handler.ConfirmMedia)
		}
	}

	// Internal S2S routes
	internal := base.Group("/internal")
	{
		internal.GET("/showtimes/:showtimeId/base-price", config.Handler.GetShowtimeBasePrice)
	}

	addr := fmt.Sprintf(":%s", config.Port)
	return r.Run(addr)
}
