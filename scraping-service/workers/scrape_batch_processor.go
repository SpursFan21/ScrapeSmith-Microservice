//ScrapeSmith\scraping-service\workers\scrape_batch_processor.go

package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"scraping-service/models"
	"scraping-service/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ProcessScrapeBatchQueue() {
	coll := utils.GetCollection("queued_scrape_jobs")

	// Fetch up to 3 pending jobs
	cursor, err := coll.Find(context.TODO(), bson.M{
		"status": "pending",
	}, options.Find().SetLimit(3).SetSort(bson.M{"created_at": 1}))
	if err != nil {
		log.Println("Failed to fetch scrape jobs:", err)
		return
	}
	defer cursor.Close(context.TODO())

	var jobs []models.QueuedScrapeJob
	if err = cursor.All(context.TODO(), &jobs); err != nil {
		log.Println("Failed to decode scrape jobs:", err)
		return
	}

	if len(jobs) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(j models.QueuedScrapeJob) {
			defer wg.Done()
			processSingleScrapeJob(j)
		}(job)

		// Lock job immediately (status = processing)
		_, err := coll.UpdateOne(context.TODO(), bson.M{"order_id": job.OrderID}, bson.M{
			"$set": bson.M{
				"status":        "processing",
				"last_tried_at": time.Now(),
			},
			"$inc": bson.M{"attempts": 1},
		})
		if err != nil {
			log.Printf("Failed to lock job %s: %v", job.OrderID, err)
		}
	}

	wg.Wait()
}

func processSingleScrapeJob(job models.QueuedScrapeJob) {
	log.Printf("⚙️ Scraping job %s", job.OrderID)

	body, err := utils.ScrapeWithScrapeNinja(job.URL)
	if err != nil {
		log.Printf("Scrape failed for %s: %v", job.OrderID, err)
		retryOrFail(job.OrderID, job.Attempts)
		return
	}

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
		log.Printf("Failed to store scrape result for %s: %v", job.OrderID, err)
		failJob(job.OrderID)
		return
	}

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
		log.Printf("Failed to send to cleaner for job %s: %v", job.OrderID, err)
		failJob(job.OrderID)
		return
	}
	defer resp.Body.Close()
	log.Printf("Job %s cleaned successfully", job.OrderID)

	// Mark as done
	utils.GetCollection("queued_scrape_jobs").UpdateOne(context.TODO(),
		bson.M{"order_id": job.OrderID},
		bson.M{"$set": bson.M{"status": "done"}},
	)
}

func retryOrFail(orderID string, attempts int) {
	status := "failed"
	if attempts < 3 {
		status = "pending"
	}
	utils.GetCollection("queued_scrape_jobs").UpdateOne(context.TODO(),
		bson.M{"order_id": orderID},
		bson.M{"$set": bson.M{"status": status}},
	)
}

func failJob(orderID string) {
	utils.GetCollection("queued_scrape_jobs").UpdateOne(context.TODO(),
		bson.M{"order_id": orderID},
		bson.M{"$set": bson.M{"status": "failed"}},
	)
}
