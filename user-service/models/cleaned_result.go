// models/cleaned_result.go

package models

import "time"

type CleanedResult struct {
	OrderID      string    `json:"order_id" bson:"orderId"`
	UserID       string    `json:"user_id" bson:"userId"`
	CreatedAt    time.Time `json:"created_at" bson:"createdAt"`
	URL          string    `json:"url" bson:"url"`
	AnalysisType string    `json:"analysis_type" bson:"analysisType"`
	CustomScript string    `json:"custom_script,omitempty" bson:"customScript,omitempty"`
	CleanedData  string    `json:"cleaned_data" bson:"cleanedData"`
}
