// ScrapeSmith\payment-service\routes\payment_routes.go

package routes

import (
	"payment-service/handlers"
	"payment-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupPaymentRoutes(app *fiber.App) {
	// Existing single-scrape endpoints
	app.Post("/create-payment-intent/:id", handlers.CreatePaymentIntent)
	app.Get("/product/:id", handlers.GetProductDetails)
	app.Post("/validate-voucher", handlers.ValidateVoucher)

	// Forge Balance endpoints with JWT middleware
	balanceGroup := app.Group("/balance", middleware.JWTMiddleware())
	balanceGroup.Post("/top-up/voucher", handlers.TopUpWithVoucher)
	balanceGroup.Get("/", handlers.GetUserBalance)
}
