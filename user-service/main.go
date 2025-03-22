package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver

	"user-service/database"
	"user-service/handlers"
)

// Load environment variables
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	// Connect to the database
	database.Connect()
	defer database.DB.Close()

	app := fiber.New()

	// Middleware to protect routes
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("AUTH0_CLIENT_SECRET")),
	}))

	// Routes
	app.Get("/users/:id", func(c *fiber.Ctx) error {
		return handlers.GetUser(c, database.DB)
	})
	app.Put("/users/:id", func(c *fiber.Ctx) error {
		return handlers.UpdateUser(c, database.DB)
	})

	// Start the server
	log.Println("User service running on port 3001")
	log.Fatal(app.Listen(":3001"))
}
