//ScrapeSmith\user-service\routes\user_routes.go

package routes

import (
	"user-service/handlers"
	"user-service/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/users")

	// User profile routes
	userGroup.Get("/:id", handlers.GetUser)
	userGroup.Put("/:id", handlers.UpdateUser)
	userGroup.Put("/:id/password", handlers.UpdatePassword)

	// Job and order routes
	userGroup.Get("/me/completed-jobs", handlers.GetCompletedJobs)

	// Unified order details route (replaces the 3 separate ones)
	userGroup.Get("/order-details/:orderId", handlers.GetFullOrderDetails)

	// Support ticket routes
	app.Post("/users/tickets", middleware.JWTMiddleware(), handlers.SubmitTicket)
	app.Get("/users/tickets", middleware.JWTMiddleware(), handlers.GetMyTickets)
	app.Post("/users/tickets/:id/reply", middleware.JWTMiddleware(), handlers.ReplyToTicket)

	// (Optional: Keep these if needed for admin tools or debugging)
	// userGroup.Get("/scraped-order/:orderId", handlers.GetScrapedOrderByID)
	// userGroup.Get("/cleaned-order/:orderId", handlers.GetCleanedOrderByID)
	// userGroup.Get("/orders/:orderId", handlers.GetOrderMetadata)
	// userGroup.Get("/ai-analysis/:orderId", handlers.GetAIAnalysisByOrderId)
}
