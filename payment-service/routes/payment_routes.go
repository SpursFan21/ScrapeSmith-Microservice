package routes

import (
	"payment-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupPaymentRoutes(app *fiber.App) {
	app.Post("/create-checkout-session", handlers.CreateCheckoutSession)
}
