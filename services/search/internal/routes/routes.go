package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes() error {
	r := gin.New()
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		search := v1.Group("/search")
		{
			search.GET("/health", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"status": true})
			})
		}
	}

	return r.Run(":6007")
}
