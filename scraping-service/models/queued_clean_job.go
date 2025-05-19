//ScrapeSmith\scraping-service\models\queued_clean_job.go

package models

import "time"

type QueuedCleanJob struct {
	OrderID      string    `bson:"orderId"`
	UserID       string    `bson:"userId"`
	URL          string    `bson:"url"`
	AnalysisType string    `bson:"analysisType"`
	CustomScript string    `bson:"customScript,omitempty"`
	CreatedAt    time.Time `bson:"createdAt"`
	RawData      string    `bson:"rawData"`
	Status       string    `bson:"status"`
	Attempts     int       `bson:"attempts"`
	LastTriedAt  time.Time `bson:"lastTriedAt,omitempty"`
	Error        string    `bson:"error,omitempty"`
}
