package repository

import (
	"ad-tracking-system/internal/domain/models"
	"database/sql"
	"log"
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
