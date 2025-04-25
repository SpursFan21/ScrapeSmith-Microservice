// user-service\handlers\completed_jobs.go
package handlers

import (
	"context"
	"log"
	"time"

	"user-service/models"
	"user-service/mongo"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCompletedJobs(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := *(user.Claims.(*jwt.MapClaims))

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Println("Failed to extract user ID from token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	collection := mongo.GetCollection("scraped_data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Updated to match snake_case as stored in MongoDB
	projection := bson.M{
		"order_id":      1,
		"user_id":       1,
		"created_at":    1,
		"url":           1,
		"analysis_type": 1,
		"custom_script": 1,
	}

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetProjection(projection)

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		log.Printf("MongoDB query failed: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch completed jobs",
		})
	}
	defer cursor.Close(ctx)

	var results []models.LightScrapeResult
	if err := cursor.All(ctx, &results); err != nil {
		log.Printf("Failed to decode MongoDB results: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse job data",
		})
	}

	return c.JSON(results)
}
