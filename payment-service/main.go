//ScrapeSmith\payment-service\main.go

package main

import (
	"log"
	"payment-service/routes"
	"payment-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env config
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize Stripe and MongoDB
	utils.InitStripe()
	utils.InitMongo()

	app := fiber.New()

	// Register API routes
	routes.SetupPaymentRoutes(app)

	log.Println("Payment service running on port 3002")
	log.Fatal(app.Listen(":3002"))
}
