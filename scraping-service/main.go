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

	// Start background scrape worker
	go func() {
		for {
			workers.ProcessScrapeQueue()
			time.Sleep(5 * time.Second)
		}
	}()

	log.Println("Scraping service running on port 3003")
	log.Fatal(app.Listen(":3003"))
}
