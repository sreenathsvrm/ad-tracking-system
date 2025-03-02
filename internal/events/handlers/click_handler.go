package handlers

import (
	"ad-tracking-system/internal/domain/models"
	"ad-tracking-system/internal/repository"
	"encoding/json"
	"log"
)

func HandleClickEvent(message []byte, repo *repository.ClickRepository) {
	var click models.ClickEvent
	if err := json.Unmarshal(message, &click); err != nil {
		log.Printf("Failed to unmarshal click event: %v", err)
		return
	}

	// Save the click event to the database
	if err := repo.Save(click); err != nil {
		log.Printf("Failed to save click event: %v", err)
		return
	}
}
