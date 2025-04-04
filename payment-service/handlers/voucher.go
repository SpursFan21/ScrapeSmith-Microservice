package handlers

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type VoucherRequest struct {
	Code string `json:"code"`
}

func ValidateVoucher(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid Authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server misconfiguration: missing JWT_SECRET_KEY",
		})
	}

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Parse request body
	var req VoucherRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	expectedCode := strings.TrimSpace(os.Getenv("VOUCHER_CODE"))
	submittedCode := strings.TrimSpace(req.Code)

	if submittedCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"valid":   false,
			"message": "Voucher code cannot be empty",
		})
	}

	if submittedCode == expectedCode {
		return c.JSON(fiber.Map{
			"valid":   true,
			"message": "Voucher code is valid",
		})
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"valid":   false,
		"message": "Invalid voucher code",
	})
}
