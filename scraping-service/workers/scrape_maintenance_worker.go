// ScrapeSmith\scraping-service\workers\scrape_maintenance_worker.go
package workers

import (
	"context"
	"log"
	"time"

	"scraping-service/utils"

	"go.mongodb.org/mongo-driver/bson"
)

func RunScrapeQueueMaintenance() {
	coll := utils.GetCollection("queued_scrape_jobs")

	// 1. Delete "done" jobs older than 24h
	cutoff := time.Now().Add(-24 * time.Hour)
	delDone, err := coll.DeleteMany(context.TODO(), bson.M{
		"status":     "done",
		"created_at": bson.M{"$lt": cutoff},
	})
	if err != nil {
		log.Printf("Cleanup error (done): %v", err)
	} else if delDone.DeletedCount > 0 {
		log.Printf("🧹 Cleaned %d old done jobs", delDone.DeletedCount)
	}

	// 2. Delete permanently failed jobs
	delFailed, err := coll.DeleteMany(context.TODO(), bson.M{
		"status": "permanently_failed",
	})
	if err != nil {
		log.Printf("Cleanup error (permanently_failed): %v", err)
	} else if delFailed.DeletedCount > 0 {
		log.Printf("🪦 Deleted %d permanently failed jobs", delFailed.DeletedCount)
	}

	// 3. Recover stuck "processing" jobs (older than 10m)
	stuckCutoff := time.Now().Add(-10 * time.Minute)
	rescue, err := coll.UpdateMany(context.TODO(), bson.M{
		"status":        "processing",
		"last_tried_at": bson.M{"$lt": stuckCutoff},
	}, bson.M{
		"$set": bson.M{
			"status": "failed", // Mark stuck jobs as "failed"
		},
	})
	if err != nil {
		log.Printf("Error rescuing stuck jobs: %v", err)
	} else if rescue.ModifiedCount > 0 {
		log.Printf("Recovered %d stuck jobs and marked as failed", rescue.ModifiedCount)
	}
}
