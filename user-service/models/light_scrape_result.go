// user-service/models/light_scrape_result.go

package models

import "time"

type LightScrapeResult struct {
	OrderID      string    `json:"orderId" bson:"orderId"`
	UserID       string    `json:"userId" bson:"userId"`
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	URL          string    `json:"url" bson:"url"`
	AnalysisType string    `json:"analysisType" bson:"analysisType"`
	CustomScript string    `json:"customScript,omitempty" bson:"customScript,omitempty"`
}
