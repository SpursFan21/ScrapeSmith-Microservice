package models

import "time"

type ScrapeRequest struct {
	URL          string `json:"url"`
	UserID       string `json:"user_id"`
	AnalysisType string `json:"analysis_type"`
	CustomScript string `json:"custom_script,omitempty"`
}

type ScrapeResult struct {
	OrderID      string    `json:"order_id"`
	UserID       string    `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	URL          string    `json:"url"`
	AnalysisType string    `json:"analysis_type"`
	CustomScript string    `json:"custom_script,omitempty"`
	Data         string    `json:"data"`
}
