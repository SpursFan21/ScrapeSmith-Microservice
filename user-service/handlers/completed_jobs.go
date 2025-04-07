//user-service\handlers\completed_jobs.go
package handlers

import (
	"context"
	"log"
	"time"

	"user-service/models"
	"user-service/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetCompletedJobs(c *fiber.Ctx) error {
	// Get user claims from middleware
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	collection := mongo.GetCollection("scraped_data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Filter jobs by user ID
	cursor, err := collection.Find(ctx, fiber.Map{"user_id": userID})
	if err != nil {
		log.Printf("MongoDB query failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch completed jobs",
		})
	}
	defer cursor.Close(ctx)

	var results []models.ScrapeResult
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Failed to decode MongoDB results: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse job data",
		})
	}

	return c.JSON(results)
}
