//ScrapeSmith\payment-service\handlers\get_balance.go
// forge balance

package handlers

import (
	"context"
	"time"

	"payment-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUserBalance(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := *(user.Claims.(*jwt.MapClaims))
	userID := claims["sub"].(string)

	collection := utils.GetCollection("user_balances")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result struct {
		Balance int `bson:"balance"`
	}

	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Balance not found"})
	}

	return c.JSON(fiber.Map{"balance": result.Balance})
}
