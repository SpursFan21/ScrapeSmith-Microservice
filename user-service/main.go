// user-service\main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"user-service/middleware"
	"user-service/mongo"
	"user-service/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to MongoDB
	mongo.ConnectMongo()

	app := fiber.New()

	// Apply JWT middleware
	app.Use(middleware.JWTMiddleware())

	// Setup routes (no longer pass db)
	routes.SetupUserRoutes(app)

	log.Println("✅ User service running on port 3001")
	log.Fatal(app.Listen(":3001"))
}
