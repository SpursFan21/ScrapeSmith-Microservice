//user-service\models\scrape_result.go
package models

import "time"

type ScrapeResult struct {
	OrderID      string    `json:"order_id" bson:"order_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	URL          string    `json:"url" bson:"url"`
	AnalysisType string    `json:"analysis_type" bson:"analysis_type"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}
