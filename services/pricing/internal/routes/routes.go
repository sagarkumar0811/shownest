package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/pkg/jwt"
	"github.com/shownest/pkg/middleware"
	"github.com/shownest/pricing-service/internal/handlers"
)

type Config struct {
	Handler    *handlers.Handler
	JWTService *jwt.Service
	Port       string
}

func InitRoutes(cfg Config) error {
	r := gin.New()
	r.Use(gin.Recovery())

	base := r.Group("/api/pricing")

	base.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})

	v1 := base.Group("/v1")
	{
		auth := middleware.JWTAuth(cfg.JWTService)
		merchantOnly := middleware.RequireMerchant()

		// Authenticated user routes
		user := v1.Group("", auth)
		{
			user.POST("/calculate", cfg.Handler.CalculatePrice)
			user.POST("/coupons/validate", cfg.Handler.ValidateCoupon)
		}

		// Merchant-only routes
		merchant := v1.Group("", auth, merchantOnly)
		{
			merchant.POST("/coupons", cfg.Handler.CreateCoupon)
			merchant.GET("/coupons", cfg.Handler.ListMyCoupons)
			merchant.POST("/rules", cfg.Handler.CreatePricingRule)
			merchant.GET("/rules", cfg.Handler.ListPricingRules)
		}
	}

	// Internal S2S routes
	internal := base.Group("/internal")
	{
		internal.POST("/coupons/:couponCode/redemptions", cfg.Handler.RecordRedemption)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	return r.Run(addr)
}
