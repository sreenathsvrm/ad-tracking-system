package services

import (
	"ad-tracking-system/internal/domain/models"
	"ad-tracking-system/internal/repository"
	"log"
)

// AdService provides business logic for ad-related operations
type AdService struct {
	adRepo *repository.AdRepository
}

// NewAdService creates a new AdService
func NewAdService(adRepo *repository.AdRepository) *AdService {
	return &AdService{adRepo: adRepo}
}

// GetAllAds fetches all ads from the repository
func (s *AdService) GetAllAds() ([]models.Ad, error) {
	ads, err := s.adRepo.FetchAll()
	if err != nil {
		log.Printf("Failed to fetch ads: %v", err)
		return nil, err
	}
	return ads, nil
}
