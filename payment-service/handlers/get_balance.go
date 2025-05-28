//ScrapeSmith\payment-service\handlers\get_balance.go
// forge balance

package handlers

import (
	"context"
	"time"

	"payment-service/models"
	"payment-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserBalance(c *fiber.Ctx) error {
	localUser := c.Locals("user")
	if localUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: no token provided",
		})
	}

	claims, ok := localUser.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: invalid token format",
		})
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user ID missing",
		})
	}

	collection := utils.GetCollection("user_balances")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result models.UserBalance
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)

	// If not found, insert default balance
	if err == mongo.ErrNoDocuments {
		result = models.UserBalance{
			UserID:      userID,
			Balance:     0,
			LastUpdated: time.Now(),
		}
		_, insertErr := collection.InsertOne(ctx, result)
		if insertErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to initialize balance",
			})
		}
		// Return new balance
		return c.JSON(fiber.Map{"balance": result.Balance})
	}

	// Other DB error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch balance",
		})
	}

	return c.JSON(fiber.Map{"balance": result.Balance})
}
