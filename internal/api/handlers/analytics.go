package handlers

import (
	"ad-tracking-system/internal/domain/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAnalytics returns analytics for a specific ad
func GetAnalytics(analyticsService *services.ClickService) gin.HandlerFunc {
	return func(c *gin.Context) {
		adID := c.Query("ad_id")
		if adID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ad_id is required"})
			return
		}

		// Get the click count for the ad
		count, err := analyticsService.GetClickCount(adID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"ad_id": adID, "click_count": count})
	}
}
