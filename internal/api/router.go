package api

import (
	"ad-tracking-system/internal/api/handlers"
	"ad-tracking-system/internal/domain/services"

	"github.com/gin-gonic/gin"
)

func NewRouter(adService *services.AdService) *gin.Engine {
	router := gin.Default()

	// API routes
	router.GET("/ads", func(c *gin.Context) {
		handlers.GetAds(c, adService)
	})

	return router
}
