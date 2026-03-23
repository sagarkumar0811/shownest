package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		booking := v1.Group("/booking")
		{
			booking.GET("/health", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": true})
			})
		}
	}

	return r
}
