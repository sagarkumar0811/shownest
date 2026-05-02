package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/booking-service/internal/handlers"
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

	base := r.Group("/api/booking")

	base.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})

	auth := middleware.JWTAuth(config.JWTService)

	// Authenticated user routes
	v1 := base.Group("/v1", auth)
	{
		v1.POST("/orders/checkout", config.Handler.CreateBooking)
		v1.GET("/orders", config.Handler.ListUserBookings)
		v1.GET("/orders/:id", config.Handler.GetBooking)
		v1.POST("/orders/:id/confirm", config.Handler.ConfirmBooking)
		v1.POST("/orders/:id/cancel", config.Handler.CancelBooking)
	}

	addr := fmt.Sprintf(":%s", config.Port)
	return r.Run(addr)
}
