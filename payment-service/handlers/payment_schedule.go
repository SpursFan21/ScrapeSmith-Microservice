//ScrapeSmith\payment-service\handlers\payment_schedule.go

package handlers

import (
	"context"
	"payment-service/models"
	"payment-service/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScheduleRequest struct {
	Amount int `json:"amount"`
}

func DeductForgeBalance(c *fiber.Ctx) error {
	// 1. Auth - extract user from JWT
	localUser := c.Locals("user")
	if localUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := localUser.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	userID, ok := claims["sub"].(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing user ID in token"})
	}

	// 2. Parse request body
	var req ScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}
	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Amount must be greater than 0"})
	}

	// 3. Connect to MongoDB
	collection := utils.GetCollection("user_balances")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 4. Fetch balance
	var balance models.UserBalance
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&balance)
	if err == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User balance not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user balance"})
	}

	// 5. Validate funds
	if balance.Balance < req.Amount {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Insufficient balance",
			"balance": balance.Balance,
			"required": req.Amount,
		})
	}

	// 6. Deduct and update
	newBalance := balance.Balance - req.Amount
	update := bson.M{
		"$set": bson.M{
			"balance":      newBalance,
			"last_updated": time.Now(),
		},
	}
	_, updateErr := collection.UpdateOne(ctx, bson.M{"user_id": userID}, update, options.Update())
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update balance"})
	}

	return c.JSON(fiber.Map{
		"message": "Balance deducted successfully",
		"balance": newBalance,
	})
}
