//ScrapeSmith\scraping-service\workers\scrape_retry_worker.go

package workers

import (
	"context"
	"log"
	"time"

	"scraping-service/utils"

	"go.mongodb.org/mongo-driver/bson"
)

// RetryFailedScrapeJobs finds failed jobs and resets them to pending if eligible for retry.
func RetryFailedScrapeJobs() {
	coll := utils.GetCollection("queued_scrape_jobs")

	filter := bson.M{
		"status":   "failed",
		"attempts": bson.M{"$lt": 3},
		"lastTriedAt": bson.M{
			"$lte": time.Now().Add(-30 * time.Second),
		},
	}

	update := bson.M{
		"$set": bson.M{
			"status":      "pending",
			"lastTriedAt": time.Now(),
		},
	}

	result, err := coll.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Printf("Retry worker error: %v", err)
		return
	}

	if result.ModifiedCount > 0 {
		log.Printf("Retry worker reset %d job(s) to pending", result.ModifiedCount)
	}
}
