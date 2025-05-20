// user-service\handlers\order_details.go

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
)

// for order details page
func GetFullOrderDetails(c *fiber.Ctx) error {
	orderId := c.Params("orderId")
	user := c.Locals("user").(*jwt.Token)
	claims := *(user.Claims.(*jwt.MapClaims))
	requesterID := claims["sub"].(string)
	isAdmin := claims["is_admin"].(bool)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Step 1: Fetch raw scrape data
	var raw models.ScrapeResult
	if err := mongo.GetCollection("scraped_data").FindOne(ctx, bson.M{"orderId": orderId}).Decode(&raw); err != nil {
		log.Printf("Raw data not found for %s", orderId)
		return fiber.NewError(fiber.StatusNotFound, "Raw scrape data not found")
	}

	// Access check
	if !isAdmin && raw.UserID != requesterID {
		return fiber.NewError(fiber.StatusForbidden, "Unauthorized to view this order")
	}

	// Step 2: Fetch cleaned data
	var clean models.CleanedResult
	if err := mongo.GetCollection("cleaned_data").FindOne(ctx, bson.M{"orderId": orderId}).Decode(&clean); err != nil {
		log.Printf("Cleaned data not found for %s", orderId)
		// Not fatal — continue
	}

	// Step 3: Fetch AI analysis
	var ai models.AIAnalysisResult
	if err := mongo.GetCollection("analyzed_data").FindOne(ctx, bson.M{"orderId": orderId}).Decode(&ai); err != nil {
		log.Printf("AI analysis not found for %s", orderId)
		// Not fatal — continue
	}

	// Unified response
	return c.JSON(fiber.Map{
		"order": fiber.Map{
			"order_id":      raw.OrderID,
			"created_at":    raw.CreatedAt,
			"url":           raw.URL,
			"analysis_type": raw.AnalysisType,
			"custom_script": raw.CustomScript,
		},
		"raw_data":    raw.Data,
		"clean_data":  clean.CleanedData,
		"ai_analysis": ai.AnalysisData,
	})
}
