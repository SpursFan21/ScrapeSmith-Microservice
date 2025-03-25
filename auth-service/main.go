package main

import (
	"auth-service/config"
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/middlewares"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func initDB() {
	database.Connect()
}

func main() {
	config.LoadEnv()
	initDB()
	handlers.SetDB(database.DB)

	app := fiber.New()

	// Public route (No Auth Required)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Auth Service"})
	})

	// Auth routes
	app.Post("/signup", handlers.Signup)
	app.Post("/login", handlers.Login)
	app.Post("/refresh", handlers.RefreshToken)

	// Protected routes (using JWT middleware)
	app.Use(middlewares.JWTMiddleware())
	app.Get("/user/profile", handlers.GetUserProfile)
	app.Put("/user/profile", handlers.UpdateUserProfile)

	log.Println("✅ Auth service running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
