// ScrapeSmith\payment-service\handlers\top_up.go
// forge balance

package handlers

import (
	"context"
	"os"
	"strings"
	"time"

	"payment-service/models"
	"payment-service/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TopUpVoucherRequest struct {
	Code   string `json:"code"`
	Amount int    `json:"amount"` // new field to indicate desired top-up amount (jobs)
}

/*var discountMap = map[int]int{
	10:    10,   // $1.00/job
	25:    24,   // $0.96/job
	50:    45,   // $0.90/job
	100:   85,   // $0.85/job
	500:   400,  // $0.80/job
	1000:  750,  // $0.75/job
	10000: 4000, // $0.40/job
}*/

func TopUpWithVoucher(c *fiber.Ctx) error {
	// 1. Authenticate user
	localUser := c.Locals("user")
	if localUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := localUser.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing user ID in token"})
	}

	// 2. Parse request
	var req TopUpVoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	expectedCode := strings.TrimSpace(os.Getenv("VOUCHER_CODE"))
	if req.Code != expectedCode {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid voucher code"})
	}

	// 3. Validate amount and price mapping
	//price, ok := discountMap[req.Amount]
	//if !ok {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unsupported top-up amount"})
	//}

	// 4. Connect to MongoDB
	collection := utils.GetCollection("user_balances")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5. Check if user already has a balance record
	var balance models.UserBalance
	err := collection.FindOne(ctx, bson.M{"user_id": userID}).Decode(&balance)
	if err != nil {
		// If not found, create a new record
		balance = models.UserBalance{
			UserID:      userID,
			Balance:     req.Amount,
			LastUpdated: time.Now(),
		}
		_, insertErr := collection.InsertOne(ctx, balance)
		if insertErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user balance"})
		}
	} else {
		// If found, update balance
		newBalance := balance.Balance + req.Amount
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
		balance.Balance = newBalance
	}

	return c.JSON(fiber.Map{
		"message": "Top-up successful",
		"balance": balance.Balance,
	})
}
