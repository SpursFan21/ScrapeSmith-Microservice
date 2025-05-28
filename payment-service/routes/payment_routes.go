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

	// Forge Balance endpoints (use explicit routes instead of a group)
	app.Get("/balance", middleware.JWTMiddleware(), handlers.GetUserBalance)
	app.Post("/balance/top-up/voucher", middleware.JWTMiddleware(), handlers.TopUpWithVoucher)

	// Temporary: no JWT
	//app.Get("/balance", handlers.GetUserBalance)
	//app.Post("/balance/top-up/voucher", handlers.TopUpWithVoucher)

	app.Get("/debug", func(c *fiber.Ctx) error {
		return c.SendString(" Payment service is reachable")
	})

}
