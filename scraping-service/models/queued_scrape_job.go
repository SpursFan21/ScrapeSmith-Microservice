//ScrapeSmith\scraping-service\models\queued_scrape_job.go
package models

import "time"

type QueuedScrapeJob struct {
	OrderID      string    `bson:"order_id"`
	UserID       string    `bson:"user_id"`
	CreatedAt    time.Time `bson:"created_at"`
	URL          string    `bson:"url"`
	AnalysisType string    `bson:"analysis_type"`
	CustomScript string    `bson:"custom_script,omitempty"`
	Status       string    `bson:"status"` // "pending", "processing", "done", "failed"
	Attempts     int       `bson:"attempts"`
	LastTriedAt  time.Time `bson:"last_tried_at,omitempty"`
}
