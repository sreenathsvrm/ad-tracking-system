package handlers

import (
	"ad-tracking-system/internal/domain/models"
	"ad-tracking-system/internal/domain/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RecordClick records a click event
func RecordClick(c *gin.Context, clickService *services.ClickService) {
	var click models.ClickEvent
	if err := c.ShouldBindJSON(&click); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Set the timestamp to the current time
	click.Timestamp = time.Now()
	click.IP = c.ClientIP()

	// Record the click event
	if err := clickService.RecordClick(click); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record click"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Click recorded"})
}
