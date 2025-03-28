package routes

import (
	"payment-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupPaymentRoutes(app *fiber.App) {
	app.Post("/create-payment-intent/:id", handlers.CreatePaymentIntent)
	app.Get("/product/:id", handlers.GetProductDetails)
}
