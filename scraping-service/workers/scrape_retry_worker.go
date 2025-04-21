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

	// Find 1–2 failed jobs eligible for retry (30s+ since last attempt, max 3 attempts)
	filter := bson.M{
		"status":   "failed",
		"attempts": bson.M{"$lt": 3},
		"last_tried_at": bson.M{
			"$lte": time.Now().Add(-30 * time.Second),
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":        "pending",
			"last_tried_at": time.Now(),
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
