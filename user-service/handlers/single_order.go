// user-service\handlers\single_order.go
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

func GetScrapedOrderByID(c *fiber.Ctx) error {
	orderId := c.Params("orderId")
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	requesterID := (*claims)["sub"].(string)
	isAdmin := (*claims)["is_admin"].(bool)

	collection := mongo.GetCollection("scraped_data")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.ScrapeResult
	err := collection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&result)
	if err != nil {
		log.Printf("Failed to find scraped order: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Scraped order not found",
		})
	}

	if !isAdmin && result.UserID != requesterID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized access to this order",
		})
	}

	return c.JSON(result)
}

func GetCleanedOrderByID(c *fiber.Ctx) error {
	orderId := c.Params("orderId")
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	requesterID := (*claims)["sub"].(string)
	isAdmin := (*claims)["is_admin"].(bool)

	collection := mongo.GetCollection("cleaneddatas")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.CleanedResult

	log.Printf("🔍 Looking for cleaned order with orderId = %s", orderId)

	err := collection.FindOne(ctx, bson.M{"orderId": orderId}).Decode(&result)

	if err != nil {
		log.Printf("Failed to find cleaned order: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cleaned order not found",
		})
	}

	if !isAdmin && result.UserID != requesterID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized access to this order",
		})
	}

	return c.JSON(result)
}
