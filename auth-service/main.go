// ScrapeSmith\auth-service\main.go
package main

import (
	"auth-service/config"
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load env vars
	config.LoadEnv()

	// Connect to MongoDB
	database.Connect()

	// Start Fiber app
	app := fiber.New()

	// Public route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Auth Service"})
	})

	// Auth routes
	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)
	app.Post("/refresh", handlers.RefreshToken)
	app.Post("/logout", handlers.Logout)

	// Protected routes
	app.Use(middlewares.JWTMiddleware())
	app.Get("/user/profile", handlers.GetUserProfile)
	app.Put("/user/profile", handlers.UpdateUserProfile)

	log.Println("✅ Auth service running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
