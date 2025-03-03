package repository

import (
	"ad-tracking-system/internal/domain/models"
	"database/sql"
	"log"
	"net"
)

// ClickRepository manages database operations for click events
type ClickRepository struct {
	db *sql.DB
}

// NewClickRepository creates a new ClickRepository
func NewClickRepository(db *sql.DB) *ClickRepository {
	return &ClickRepository{db: db}
}

// Save saves a click event to the database
func (r *ClickRepository) Save(click models.ClickEvent) error {
	query := `INSERT INTO clicks (ad_id, timestamp, ip, playback_time) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, click.AdID, click.Timestamp, click.IP, click.PlaybackTime)
	if err != nil {
		log.Printf("Failed to save click event: %v", err)
		return err
	}
	return nil
}

// AdExists checks if an ad with the given ID exists
func (r *ClickRepository) AdExists(adID string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM ads WHERE id = $1)"
	err := r.db.QueryRow(query, adID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// GetClickCountByIP returns the number of clicks from a specific IP in the last hour
func (r *ClickRepository) GetClickCountByIP(ip string) (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM clicks WHERE ip = $1 AND timestamp > NOW() - INTERVAL '1 hour'"
	err := r.db.QueryRow(query, ip).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// IsValidIP checks if the IP address is valid
func (r *ClickRepository) IsValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

// IsPlaybackTimeValid checks if the playback time is within a valid range
func (r *ClickRepository) IsPlaybackTimeValid(playbackTime int) bool {
	return playbackTime >= 0 && playbackTime <= 3600 // Example: 0 to 3600 seconds (1 hour)
}
