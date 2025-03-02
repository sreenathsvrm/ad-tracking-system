package repository

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

// AnalyticsRepository manages analytics data in Redis
type AnalyticsRepository struct {
	redisClient *redis.Client
}

// NewAnalyticsRepository creates a new AnalyticsRepository
func NewAnalyticsRepository(redisClient *redis.Client) *AnalyticsRepository {
	return &AnalyticsRepository{redisClient: redisClient}
}

// IncrementClickCount increments the click count for a specific ad
func (r *AnalyticsRepository) IncrementClickCount(adID string) error {
	ctx := context.Background()
	key := "clicks:" + adID
	if err := r.redisClient.Incr(ctx, key).Err(); err != nil {
		log.Printf("Failed to increment click count: %v", err)
		return err
	}
	return nil
}

// GetClickCount returns the total click count for a specific ad
func (r *AnalyticsRepository) GetClickCount(adID string) (int64, error) {
	ctx := context.Background()
	key := "clicks:" + adID
	count, err := r.redisClient.Get(ctx, key).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil // No clicks recorded yet
		}
		log.Printf("Failed to get click count: %v", err)
		return 0, err
	}
	return count, nil
}
