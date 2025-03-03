package repository

import (
	"ad-tracking-system/internal/domain/models"
	"database/sql"
	"log"
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

// CountAds returns the number of ads in the database
func (r *AdRepository) CountAds() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM ads`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Failed to count ads: %v", err)
		return 0, err
	}
	return count, nil
}

// Seed inserts 10 dummy ads into the ads table
func (r *AdRepository) Seed() error {
	// Define 10 dummy ads
	dummyAds := []models.Ad{
		{
			ID:        "1",
			ImageURL:  "https://example.com/images/ad1.jpg",
			TargetURL: "https://example.com/landing/ad1",
		},
		{
			ID:        "2",
			ImageURL:  "https://example.com/images/ad2.jpg",
			TargetURL: "https://example.com/landing/ad2",
		},
		{
			ID:        "3",
			ImageURL:  "https://example.com/images/ad3.jpg",
			TargetURL: "https://example.com/landing/ad3",
		},
		{
			ID:        "4",
			ImageURL:  "https://example.com/images/ad4.jpg",
			TargetURL: "https://example.com/landing/ad4",
		},
		{
			ID:        "5",
			ImageURL:  "https://example.com/images/ad5.jpg",
			TargetURL: "https://example.com/landing/ad5",
		},
		{
			ID:        "6",
			ImageURL:  "https://example.com/images/ad6.jpg",
			TargetURL: "https://example.com/landing/ad6",
		},
		{
			ID:        "7",
			ImageURL:  "https://example.com/images/ad7.jpg",
			TargetURL: "https://example.com/landing/ad7",
		},
		{
			ID:        "8",
			ImageURL:  "https://example.com/images/ad8.jpg",
			TargetURL: "https://example.com/landing/ad8",
		},
		{
			ID:        "9",
			ImageURL:  "https://example.com/images/ad9.jpg",
			TargetURL: "https://example.com/landing/ad9",
		},
		{
			ID:        "10",
			ImageURL:  "https://example.com/images/ad10.jpg",
			TargetURL: "https://example.com/landing/ad10",
		},
	}

	// Insert dummy ads into the database
	for _, ad := range dummyAds {
		query := `INSERT INTO ads (id, image_url, target_url) VALUES ($1, $2, $3)`
		_, err := r.db.Exec(query, ad.ID, ad.ImageURL, ad.TargetURL)
		if err != nil {
			log.Printf("Failed to insert ad %s: %v", ad.ID, err)
			return err
		}
	}

	log.Println("Successfully seeded 10 dummy ads")
	return nil
}