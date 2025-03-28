package handlers

import (
	"log"

	"scraping-service/utils"

	"github.com/gofiber/fiber/v2"
)

type ScrapeRequest struct {
	URL    string `json:"url"`
	UserID string `json:"user_id"`
}

func SingleScrape(c *fiber.Ctx) error {
	var req ScrapeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	// Perform scraping using ScrapeNinja
	rawData, err := utils.ScrapeWithScrapeNinja(req.URL)
	if err != nil {
		log.Printf("Error scraping data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to scrape data",
		})
	}

	// Store data in MongoDB
	collection := utils.GetCollection("scraped_data")
	_, err = collection.InsertOne(c.Context(), map[string]interface{}{
		"user_id": req.UserID,
		"url":     req.URL,
		"data":    string(rawData),
	})
	if err != nil {
		log.Printf("MongoDB insert error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save scraped data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Scraping successful",
	})
}
