package main

import (
	"log"
	"payment-service/routes"
	"payment-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	utils.InitStripe()

	app := fiber.New()

	routes.SetupPaymentRoutes(app)

	log.Println("Payment service running on port 3002")
	log.Fatal(app.Listen(":3002"))
}
