//ScrapeSmith\user-service\handlers\ai_analysis.go

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

func GetAIAnalysisByOrderId(c *fiber.Ctx) error {
	orderId := c.Params("orderId")

	user := c.Locals("user").(*jwt.Token)
	claims := *(user.Claims.(*jwt.MapClaims))
	userID, ok := claims["sub"].(string)
	isAdmin, _ := claims["is_admin"].(bool)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	collection := mongo.GetCollection("ai_analysis_data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.AIAnalysisResult
	err := collection.FindOne(ctx, bson.M{"orderId": orderId}).Decode(&result)
	if err != nil {
		log.Printf("Failed to find AI analysis result: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "AI Analysis not found",
		})
	}

	if result.UserID != userID && !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	return c.JSON(result)
}
