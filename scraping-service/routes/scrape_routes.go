// scraping-service/routes/scrape_routes.go
package routes

import (
	"scraping-service/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupScrapeRoutes(app *fiber.App) {
	app.Post("/single", handlers.SingleScrape)
}
