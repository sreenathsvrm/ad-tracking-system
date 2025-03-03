package services

import (
	"ad-tracking-system/internal/domain/models"
	"ad-tracking-system/internal/repository"
	"ad-tracking-system/internal/utils/circuitbreaker"
	"log"

	"github.com/sony/gobreaker"
)

type AdService struct {
	adRepo *repository.AdRepository
	cb     *gobreaker.CircuitBreaker
}

func NewAdService(adRepo *repository.AdRepository) *AdService {
	return &AdService{
		adRepo: adRepo,
		cb:     circuitbreaker.NewCircuitBreaker("ad-service"), // Initialize circuit breaker
	}
}

func (s *AdService) GetAllAds() ([]models.Ad, error) {
	// Wrap database operation with circuit breaker
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.adRepo.FetchAll()
	})
	if err != nil {
		log.Printf("Failed to fetch ads (circuit breaker): %v", err)
		return nil, err
	}

	return result.([]models.Ad), nil
}
