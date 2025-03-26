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
}
