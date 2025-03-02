package services

import (
	"ad-tracking-system/internal/domain/models"
	"ad-tracking-system/internal/repository"
	"log"
)

// ClickService provides business logic for click-related operations
type ClickService struct {
	clickRepo     *repository.ClickRepository
	analyticsRepo *repository.AnalyticsRepository
}

// NewClickService creates a new ClickService
func NewClickService(clickRepo *repository.ClickRepository, analyticsRepo *repository.AnalyticsRepository) *ClickService {
	return &ClickService{
		clickRepo:     clickRepo,
		analyticsRepo: analyticsRepo,
	}
}

// RecordClick records a click event and updates analytics
func (s *ClickService) RecordClick(click models.ClickEvent) error {
	// Save the click event to the database
	if err := s.clickRepo.Save(click); err != nil {
		log.Printf("Failed to save click event: %v", err)
		return err
	}

	// Update analytics (e.g., increment click count in Redis)
	if err := s.analyticsRepo.IncrementClickCount(click.AdID); err != nil {
		log.Printf("Failed to update analytics: %v", err)
		return err
	}

	return nil
}

// GetClickCount returns the total click count for a specific ad
func (s *ClickService) GetClickCount(adID string) (int64, error) {
	count, err := s.analyticsRepo.GetClickCount(adID)
	if err != nil {
		log.Printf("Failed to get click count: %v", err)
		return 0, err
	}
	return count, nil
}
