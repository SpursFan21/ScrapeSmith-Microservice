// user-service\main.go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"user-service/database"
	"user-service/middleware"
	"user-service/mongo"
	"user-service/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	database.Connect()
	defer database.DB.Close()

	mongo.ConnectMongo()

	app := fiber.New()

	app.Use(middleware.JWTMiddleware())

	// Setup routes
	routes.SetupUserRoutes(app, database.DB)

	log.Println("User service running on port 3001")
	log.Fatal(app.Listen(":3001"))
}
