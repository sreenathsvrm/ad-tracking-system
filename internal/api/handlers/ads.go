package handlers

import (
	"ad-tracking-system/internal/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAds fetches all ads
func GetAds(c *gin.Context, adService *services.AdService) {
	ads, err := adService.GetAllAds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch ads"})
		return
	}

	c.JSON(http.StatusOK, ads)
}
