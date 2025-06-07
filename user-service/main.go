// user-service\main.go

package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"user-service/handlers"
	"user-service/middleware"
	"user-service/mongo"
	"user-service/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Connect to MongoDB
	mongo.ConnectMongo()

	// Initialize ticket collection
	handlers.InitTicketCollection(mongo.MongoClient, mongo.GetCollection("tickets").Database().Name())

	app := fiber.New()

	// Global JWT middleware
	app.Use(middleware.JWTMiddleware())

	// Routes
	routes.SetupUserRoutes(app)

	log.Println("User service running on port 3001")
	log.Fatal(app.Listen(":3001"))
}
