// user-service\models\light_scrape_result.go
package models

import "time"

type LightScrapeResult struct {
	OrderID      string    `json:"order_id" bson:"order_id"`
	UserID       string    `json:"user_id" bson:"user_id"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	URL          string    `json:"url" bson:"url"`
	AnalysisType string    `json:"analysis_type" bson:"analysis_type"`
	CustomScript string    `json:"custom_script,omitempty" bson:"custom_script,omitempty"`
}
