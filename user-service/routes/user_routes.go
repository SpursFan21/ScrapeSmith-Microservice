//ScrapeSmith\user-service\routes\user_routes.go

package routes

import (
	"user-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	// Everything is imposed with global middleware in main
	userGroup := app.Group("/users")

	// User profile routes
	userGroup.Get("/:id", handlers.GetUser)
	userGroup.Put("/:id", handlers.UpdateUser)
	userGroup.Put("/:id/password", handlers.UpdatePassword)

	// Job and order routes
	userGroup.Get("/me/completed-jobs", handlers.GetCompletedJobs)

	// Unified order details route
	userGroup.Get("/order-details/:orderId", handlers.GetFullOrderDetails)

	// ticket
	userGroup.Get("/tickets", handlers.GetMyTickets)
	userGroup.Post("/tickets", handlers.SubmitTicket)
	userGroup.Post("/tickets/:id/reply", handlers.ReplyToTicket)

}
