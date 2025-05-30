//ScrapeSmith\payment-service\handlers\payment_schedule.go

// ScrapeSmith\payment-service\handlers\payment_schedule.go

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

	// 4. Atomic deduction with filter check
	filter := bson.M{
		"user_id": userID,
		"balance": bson.M{"$gte": req.Amount},
	}
	update := bson.M{
		"$inc": bson.M{"balance": -req.Amount},
		"$set": bson.M{"last_updated": time.Now()},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update balance"})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient balance or user not found",
		})
	}

	// 5. Return success
	var updated models.UserBalance
	err = collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&updated)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Balance deducted, but failed to fetch updated balance",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Balance deducted successfully",
		"balance": updated.Balance,
	})
}
