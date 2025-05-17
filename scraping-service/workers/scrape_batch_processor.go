//ScrapeSmith\scraping-service\workers\scrape_batch_processor.go

package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"scraping-service/models"
	"scraping-service/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ProcessScrapeBatchQueue() {
	// Get the "queued_scrape_jobs" collection
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

	// Use WaitGroup to handle concurrent processing
	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(j models.QueuedScrapeJob) {
			defer wg.Done()
			processSingleScrapeJob(j)
		}(job)

		// Lock the job immediately (status = processing)
		_, err := coll.UpdateOne(context.TODO(), bson.M{"order_id": job.OrderID, "status": "pending"}, bson.M{
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

	// Wait for all workers to finish
	wg.Wait()
}

func processSingleScrapeJob(job models.QueuedScrapeJob) {
	log.Printf("Scraping job %s", job.OrderID)

	// Scrape the content using the ScrapeNinja API
	body, err := utils.ScrapeWithScrapeNinja(job.URL)
	if err != nil {
		log.Printf("Scrape failed for %s: %v", job.OrderID, err)
		failJob(job.OrderID)
		return
	}

	// Store the raw scrape data in the "scraped_data" collection
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

	// Insert the scrape job into the data-cleaning queue (MongoDB)
	cleanerColl := utils.GetCollection("queued_clean_jobs") // Target collection in data-cleaning service MongoDB
	_, err = cleanerColl.InsertOne(context.TODO(), bson.M{
		"order_id":      job.OrderID,
		"user_id":       job.UserID,
		"url":           job.URL,
		"analysis_type": job.AnalysisType,
		"custom_script": job.CustomScript,
		"raw_data":      body, // The raw data to be cleaned
		"status":        "pending",
		"created_at":    job.CreatedAt,
	})
	if err != nil {
		log.Printf("Failed to insert job into clean queue for %s: %v", job.OrderID, err)
		failJob(job.OrderID)
		return
	}

	log.Printf("Job %s added to clean queue", job.OrderID)

	// Mark original scrape job as done
	utils.GetCollection("queued_scrape_jobs").UpdateOne(context.TODO(),
		bson.M{"order_id": job.OrderID},
		bson.M{"$set": bson.M{"status": "done"}})
}

func failJob(orderID string) {
	// Update the job status to "failed" in the queue
	_, err := utils.GetCollection("queued_scrape_jobs").UpdateOne(context.TODO(),
		bson.M{"order_id": orderID},
		bson.M{"$set": bson.M{"status": "failed"}})
	if err != nil {
		log.Printf("Failed to mark job %s as failed: %v", orderID, err)
	}
}
