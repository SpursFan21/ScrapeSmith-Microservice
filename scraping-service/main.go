// scraping-service\main.go
package main

import (
	"log"
	"scraping-service/routes"
	"scraping-service/utils"
	"scraping-service/workers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	utils.ConnectMongo()

	app := fiber.New()
	routes.SetupScrapeRoutes(app)

	// Poll every 3 seconds
	go func() {
		for {
			workers.ProcessScrapeBatchQueue()
			time.Sleep(3 * time.Second)
		}
	}()

	// Retry worker every 30 seconds
	go func() {
		for {
			workers.RetryFailedScrapeJobs()
			time.Sleep(30 * time.Second)
		}
	}()

	// Queue maintenance every 30 minutes
	go func() {
		for {
			workers.RunScrapeQueueMaintenance()
			time.Sleep(30 * time.Minute)
		}
	}()

	log.Println("Scraping service running on port 3003")
	log.Fatal(app.Listen(":3003"))
}
