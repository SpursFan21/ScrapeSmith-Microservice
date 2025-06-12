// ScrapeSmith\scraping-service\workers\scrape_batch_processor.go

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

// ProcessScrapeBatchQueue checks the scrape queue for up to 3 pending jobs.
// If found, it locks them for processing and dispatches them to worker goroutines.
func ProcessScrapeBatchQueue() {
	// Get the "queued_scrape_jobs" collection
	coll := utils.GetCollection("queued_scrape_jobs")

	// Fetch up to 3 pending jobs (sorted by createdAt)
	cursor, err := coll.Find(context.TODO(), bson.M{
		"status": "pending",
	}, options.Find().SetLimit(3).SetSort(bson.M{"createdAt": 1}))
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
		_, err := coll.UpdateOne(context.TODO(), bson.M{"orderId": job.OrderID, "status": "pending"}, bson.M{
			"$set": bson.M{
				"status":      "processing",
				"lastTriedAt": time.Now(),
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

// processSingleScrapeJob handles a single job: performs scraping,
// stores the raw data, and sends it to the data-cleaning queue.
func processSingleScrapeJob(job models.QueuedScrapeJob) {
	log.Printf("Scraping job %s", job.OrderID)

	// Scrape the content using the ScrapeNinja API
	body, err := utils.ScrapeWithScrapeNinja(job.URL)
	if err != nil {
		log.Printf("Scrape failed for %s: %v", job.OrderID, err)
		failJob(job.OrderID)
		return
	}

	log.Printf("Successfully scraped job %s", job.OrderID)

	// Store the raw scrape data in the "scraped_data" collection
	scrapedColl := utils.GetCollection("scraped_data")
	_, err = scrapedColl.InsertOne(context.TODO(), bson.M{
		"orderId":      job.OrderID,
		"userId":       job.UserID,
		"createdAt":    job.CreatedAt,
		"url":          job.URL,
		"analysisType": job.AnalysisType,
		"customScript": job.CustomScript,
		"data":         body,
	})
	if err != nil {
		log.Printf("Failed to store scrape result for %s: %v", job.OrderID, err)
		failJob(job.OrderID)
		return
	}

	// Insert the job into the data-cleaning queue
	cleanerColl := utils.GetCollection("queued_clean_jobs")
	cleanJob := models.QueuedCleanJob{
		OrderID:      job.OrderID,
		UserID:       job.UserID,
		URL:          job.URL,
		AnalysisType: job.AnalysisType,
		CustomScript: job.CustomScript,
		CreatedAt:    job.CreatedAt,
		RawData:      body,
		Status:       "pending",
		Attempts:     0,
	}
	_, err = cleanerColl.InsertOne(context.TODO(), cleanJob)
	if err != nil {
		log.Printf("Failed to insert job into clean queue for %s: %v", job.OrderID, err)
		failJob(job.OrderID)
		return
	}

	log.Printf("Job %s inserted into clean queue", job.OrderID)

	// Mark original scrape job as done
	_, err = utils.GetCollection("queued_scrape_jobs").UpdateOne(context.TODO(),
		bson.M{"orderId": job.OrderID},
		bson.M{"$set": bson.M{"status": "done"}})
	if err != nil {
		log.Printf("Failed to mark job %s as done", job.OrderID)
	} else {
		log.Printf("Job %s marked as done in scrape queue", job.OrderID)
	}
}

// failJob updates the job status to "failed" in the scrape queue
func failJob(orderID string) {
	_, err := utils.GetCollection("queued_scrape_jobs").UpdateOne(context.TODO(),
		bson.M{"orderId": orderID},
		bson.M{"$set": bson.M{"status": "failed"}})
	if err != nil {
		log.Printf("Failed to mark job %s as failed: %v", orderID, err)
	}
}
