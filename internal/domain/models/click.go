package models

import "time"

// ClickEvent represents a user click on an ad
type ClickEvent struct {
	AdID         string    `json:"ad_id"`
	Timestamp    time.Time `json:"timestamp"`
	IP           string    `json:"ip"`
	PlaybackTime int       `json:"playback_time"`
}
