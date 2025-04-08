// user-service\routes\user_routes.go
package routes

import (
	"database/sql"
	"user-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(app *fiber.App, db *sql.DB) {
	userGroup := app.Group("/users")
	userGroup.Get("/:id", func(c *fiber.Ctx) error {
		return handlers.GetUser(c, db)
	})
	userGroup.Put("/:id", func(c *fiber.Ctx) error {
		return handlers.UpdateUser(c, db)
	})
	userGroup.Put("/:id/password", func(c *fiber.Ctx) error {
		return handlers.UpdatePassword(c, db)
	})

	userGroup.Get("/me/completed-jobs", handlers.GetCompletedJobs)
	userGroup.Get("/scraped-order/:orderId", handlers.GetScrapedOrderByID)
	userGroup.Get("/cleaned-order/:orderId", handlers.GetCleanedOrderByID)
	userGroup.Get("/orders/:orderId", handlers.GetOrderMetadata)

}
