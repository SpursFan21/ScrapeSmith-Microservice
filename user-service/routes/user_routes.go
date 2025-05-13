// user-service\routes\user_routes.go
// user-service\routes\user_routes.go
package routes

import (
	"user-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App) {
	userGroup := app.Group("/users")

	userGroup.Get("/:id", handlers.GetUser)
	userGroup.Put("/:id", handlers.UpdateUser)
	userGroup.Put("/:id/password", handlers.UpdatePassword)

	userGroup.Get("/me/completed-jobs", handlers.GetCompletedJobs)
	userGroup.Get("/scraped-order/:orderId", handlers.GetScrapedOrderByID)
	userGroup.Get("/cleaned-order/:orderId", handlers.GetCleanedOrderByID)
	userGroup.Get("/orders/:orderId", handlers.GetOrderMetadata)
	userGroup.Get("/ai-analysis/:orderId", handlers.GetAIAnalysisByOrderId)
}
