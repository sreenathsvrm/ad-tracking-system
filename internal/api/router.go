package api

import (
	"ad-tracking-system/internal/api/handlers"
	"ad-tracking-system/internal/domain/services"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewRouter initializes the API routes and middleware
func NewRouter(adService *services.AdService, clickService *services.ClickService) *gin.Engine {
	router := gin.Default()

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API routes
	router.GET("/ads", func(c *gin.Context) {
		handlers.GetAds(c, adService)
	})
	router.POST("/ads/click", func(c *gin.Context) {
		handlers.RecordClick(c, clickService)
	})
	router.GET("/ads/analytics", handlers.GetAnalytics(clickService))

	return router
}
