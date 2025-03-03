package services

import (
	"ad-tracking-system/internal/domain/models"
	"ad-tracking-system/internal/repository"
	"ad-tracking-system/internal/utils/circuitbreaker"
	"log"

	"github.com/sony/gobreaker"
)

type ClickService struct {
	clickRepo     *repository.ClickRepository
	analyticsRepo *repository.AnalyticsRepository
	cb            *gobreaker.CircuitBreaker
}

func NewClickService(clickRepo *repository.ClickRepository, analyticsRepo *repository.AnalyticsRepository) *ClickService {
	return &ClickService{
		clickRepo:     clickRepo,
		analyticsRepo: analyticsRepo,
		cb:            circuitbreaker.NewCircuitBreaker("click-service"), // Initialize circuit breaker
	}
}

func (s *ClickService) RecordClick(click models.ClickEvent) error {
	// Wrap database operation with circuit breaker
	_, err := s.cb.Execute(func() (interface{}, error) {
		if err := s.clickRepo.Save(click); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Printf("Failed to record click (circuit breaker): %v", err)
		return err
	}

	// Wrap Redis operation with circuit breaker
	_, err = s.cb.Execute(func() (interface{}, error) {
		if err := s.analyticsRepo.IncrementClickCount(click.AdID); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Printf("Failed to update analytics (circuit breaker): %v", err)
		return err
	}

	return nil
}

func (s *ClickService) GetClickCount(adID string) (int64, error) {
	// Wrap Redis operation with circuit breaker
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.analyticsRepo.GetClickCount(adID)
	})
	if err != nil {
		log.Printf("Failed to get click count (circuit breaker): %v", err)
		return 0, err
	}

	return result.(int64), nil
}
