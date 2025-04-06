//user-service\handlers\completed.go
package handlers

import (
	"context"
	"log"
	"time"
	"user-service/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

type CompletedJob struct {
	OrderID      string    `json:"orderId" bson:"order_id"`
	URL          string    `json:"url" bson:"url"`
	AnalysisType string    `json:"analysisType" bson:"analysis_type"`
	CreatedAt    time.Time `json:"createdAt" bson:"created_at"`
}

func GetCompletedJobs(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["sub"].(string)

	collection := database.GetMongoCollection("scraped_data")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("❌ Failed to fetch jobs for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch completed jobs",
		})
	}
	defer cursor.Close(ctx)

	var jobs []CompletedJob
	for cursor.Next(ctx) {
		var job CompletedJob
		if err := cursor.Decode(&job); err == nil {
			jobs = append(jobs, job)
		}
	}

	return c.JSON(fiber.Map{
		"jobs": jobs,
	})
}
