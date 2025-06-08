//ScrapeSmith\user-service\routes\user_routes.go

package routes

import (
	"user-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/users")

	// Static routes
	userGroup.Get("/tickets", handlers.GetMyTickets)
	userGroup.Post("/tickets", handlers.SubmitTicket)
	userGroup.Post("/tickets/:id/reply", handlers.ReplyToTicket)

	userGroup.Get("/me/completed-jobs", handlers.GetCompletedJobs)
	userGroup.Get("/order-details/:orderId", handlers.GetFullOrderDetails)

	// Dynamic routes
	userGroup.Get("/:id", handlers.GetUser)
	userGroup.Put("/:id", handlers.UpdateUser)
	userGroup.Put("/:id/password", handlers.UpdatePassword)
}
