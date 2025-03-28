package models

type ScrapeRequest struct {
	URL    string `json:"url"`
	UserID string `json:"user_id"`
}

type ScrapeResult struct {
	UserID string `json:"user_id"`
	URL    string `json:"url"`
	Data   string `json:"data"`
}
