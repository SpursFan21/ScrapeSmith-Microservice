//ScrapeSmith\scraping-service\workers\scrape_retry_worker.go

package workers

import (
	"context"
	"log"
	"time"

	"scraping-service/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func RetryFailedScrapeJobs() {
	coll := utils.GetCollection("queued_scrape_jobs")

	// Find failed jobs eligible for retry (30s+ since last attempt, max 3 attempts)
	filter := bson.M{
		"status":   "failed",
		"attempts": bson.M{"$lt": 3}, // Max 3 attempts
		"lastTriedAt": bson.M{
			"$lte": time.Now().Add(-30 * time.Second), // Retry after 30s from last attempt
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":      "pending",  // Reset status to "pending" for retry
			"lastTriedAt": time.Now(), // Update last tried time
		},
	}

	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Retry worker error: %v", err)
		return
	}

	if result.ModifiedCount > 0 {
		log.Printf("Retry worker: %d job(s) reset to pending", result.ModifiedCount)
	}
}
