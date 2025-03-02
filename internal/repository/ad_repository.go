package repository

import (
	"ad-tracking-system/internal/domain/models"
	"database/sql"
)

// AdRepository manages database operations for ads
type AdRepository struct {
	db *sql.DB
}

// NewAdRepository creates a new AdRepository
func NewAdRepository(db *sql.DB) *AdRepository {
	return &AdRepository{db: db}
}

// FetchAll fetches all ads from the database
func (r *AdRepository) FetchAll() ([]models.Ad, error) {
	query := `SELECT id, image_url, target_url FROM ads`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ads []models.Ad
	for rows.Next() {
		var ad models.Ad
		if err := rows.Scan(&ad.ID, &ad.ImageURL, &ad.TargetURL); err != nil {
			return nil, err
		}
		ads = append(ads, ad)
	}
	return ads, nil
}
