// ScrapeSmith\payment-service\handlers\top_up.go
// forge balance

package handlers

import (
	"context"
	"os"
	"time"

	"payment-service/models"
	"payment-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func TopUpWithVoucher(c *fiber.Ctx) error {
	type VoucherRequest struct {
		Code string `json:"code"`
	}

	var req VoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

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

	expectedCode := os.Getenv("FORGE_VOUCHER_CODE")
	if req.Code != expectedCode {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid voucher code"})
	}

	collection := utils.GetCollection("user_balances")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var current models.UserBalance
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&current)
	if err != nil {
		// New user, create balance
		newBalance := models.UserBalance{
			UserID:      userID,
			Balance:     100,
			LastUpdated: time.Now(),
		}
		_, err := collection.InsertOne(ctx, newBalance)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to initialize balance"})
		}
		return c.JSON(fiber.Map{"message": "Voucher applied", "balance": newBalance.Balance})
	}

	// Existing user, update balance
	newAmount := current.Balance + 100
	update := bson.M{
		"$set": bson.M{
			"balance":      newAmount,
			"last_updated": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"user_id": userID}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update balance"})
	}

	return c.JSON(fiber.Map{"message": "Voucher applied", "balance": newAmount})
}
