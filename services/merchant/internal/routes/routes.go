package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shownest/merchant-service/internal/handlers"
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

	base := r.Group("/api/merchant")

	base.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})

	v1 := base.Group("/v1")
	{
		auth := middleware.JWTAuth(config.JWTService)
		merchantOnly := middleware.RequireMerchant()

		venues := v1.Group("/venues")
		{
			venues.GET("/nearby", config.Handler.GetNearbyVenues)
			venues.GET("/city/:city", config.Handler.GetVenuesByCity)
		}

		protected := v1.Group("", auth)
		{
			protected.POST("/merchants", config.Handler.CreateMerchant)
			protected.GET("/merchants/me", config.Handler.GetMyMerchant)
			protected.POST("/merchants/me/submit", config.Handler.SubmitForReview)
		}

		merchant := v1.Group("", auth, merchantOnly)
		{
			merchant.POST("/venues", config.Handler.CreateVenue)
			merchant.GET("/venues", config.Handler.ListMyVenues)
			merchant.GET("/venues/:id", config.Handler.GetVenue)
			merchant.POST("/venues/:id/halls", config.Handler.CreateHall)
			merchant.GET("/venues/:id/halls", config.Handler.ListHalls)

			merchant.POST("/merchants/me/documents/upload-url", config.Handler.RequestDocumentUploadURL)
			merchant.POST("/merchants/me/documents/confirm", config.Handler.ConfirmDocument)
			merchant.GET("/merchants/me/documents", config.Handler.ListDocuments)
		}
	}

	addr := fmt.Sprintf(":%s", config.Port)
	return r.Run(addr)
}
