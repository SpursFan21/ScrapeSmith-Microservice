package main

import (
	"log"
	"scraping-service/routes"
	"scraping-service/utils"

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

	log.Println("Scraping service running on port 3003")
	log.Fatal(app.Listen(":3003"))
}
