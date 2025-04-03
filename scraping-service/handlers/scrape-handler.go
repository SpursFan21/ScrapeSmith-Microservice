package handlers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"scraping-service/utils"
)

type ScrapeRequest struct {
	URL        string `json:"url"`
	UserID     string `json:"user_id"`
	Analysis   string `json:"analysis_type"`
	CustomCode string `json:"custom_script,omitempty"`
}

func SingleScrape(c *fiber.Ctx) error {
	var req ScrapeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	orderId := uuid.New().String()
	createdAt := time.Now()

	// Perform scraping
	rawData, err := utils.ScrapeWithScrapeNinja(req.URL)
	if err != nil {
		log.Printf("Scraping error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to scrape data",
		})
	}

	// Store in MongoDB
	collection := utils.GetCollection("scraped_data")
	_, err = collection.InsertOne(c.Context(), map[string]interface{}{
		"order_id":      orderId,
		"user_id":       req.UserID,
		"created_at":    createdAt,
		"url":           req.URL,
		"analysis_type": req.Analysis,
		"custom_script": req.CustomCode,
		"data":          rawData,
	})
	if err != nil {
		log.Printf("Mongo insert error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save scraped data",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Scraping successful",
		"order_id":   orderId,
		"created_at": createdAt,
	})
}
