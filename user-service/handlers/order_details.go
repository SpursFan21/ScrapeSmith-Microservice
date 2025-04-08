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

func GetOrderMetadata(c *fiber.Ctx) error {
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

	collection := mongo.GetCollection("scraped_data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.ScrapeResult
	err := collection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&result)
	if err != nil {
		log.Printf("❌ Failed to find scraped order metadata: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	if result.UserID != userID && !isAdmin {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Strip the full raw data for this endpoint (metadata only)
	result.Data = ""
	return c.JSON(result)
}
