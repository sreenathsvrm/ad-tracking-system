package models

// Ad represents an advertisement
type Ad struct {
	ID        string `json:"id"`
	ImageURL  string `json:"image_url"`
	TargetURL string `json:"target_url"`
}
