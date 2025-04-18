// ScrapeSmith\scraping-service\workers\scrape_processor.go
package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"scraping-service/models"
	"scraping-service/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ProcessScrapeQueue() {
	coll := utils.GetCollection("queued_scrape_jobs")

	// Find a pending job
	filter := bson.M{"status": "pending"}
	update := bson.M{
		"$set": bson.M{"status": "processing", "last_tried_at": time.Now()},
		"$inc": bson.M{"attempts": 1},
	}
	var job models.QueuedScrapeJob
	err := coll.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&job)
	if err != nil {
		log.Println("ℹNo pending scrape jobs at this time")
		return
	}

	log.Printf("⚙️  Processing scrape job %s", job.OrderID)

	// Perform scrape
	body, err := utils.ScrapeWithScrapeNinja(job.URL)
	if err != nil {
		log.Printf("Scrape failed: %v", err)
		status := "failed"
		if job.Attempts < 3 {
			status = "pending"
		}
		coll.UpdateByID(context.TODO(), job.OrderID, bson.M{"$set": bson.M{"status": status}})
		return
	}

	// Save result to scraped_data
	scrapedColl := utils.GetCollection("scraped_data")
	_, err = scrapedColl.InsertOne(context.TODO(), bson.M{
		"order_id":      job.OrderID,
		"user_id":       job.UserID,
		"created_at":    job.CreatedAt,
		"url":           job.URL,
		"analysis_type": job.AnalysisType,
		"custom_script": job.CustomScript,
		"data":          body,
	})
	if err != nil {
		log.Printf("Failed to save scraped result: %v", err)
		coll.UpdateByID(context.TODO(), job.OrderID, bson.M{"$set": bson.M{"status": "failed"}})
		return
	}

	// Send to data-cleaning-service
	payload := map[string]interface{}{
		"orderId":      job.OrderID,
		"userId":       job.UserID,
		"url":          job.URL,
		"analysisType": job.AnalysisType,
		"customScript": job.CustomScript,
		"createdAt":    job.CreatedAt,
		"rawData":      body,
	}
	jsonBody, _ := json.Marshal(payload)
	cleanerURL := os.Getenv("DATA_CLEANER_URL")
	if cleanerURL == "" {
		cleanerURL = "http://data-cleaning-service:3004/api/clean"
	}

	resp, err := http.Post(cleanerURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to send to cleaner: %v", err)
	} else {
		defer resp.Body.Close()
		log.Printf("Cleaning job dispatched: %s", resp.Status)
	}

	// Mark as done
	coll.UpdateByID(context.TODO(), job.OrderID, bson.M{"$set": bson.M{"status": "done"}})
}
