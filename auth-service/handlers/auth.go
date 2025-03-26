package handlers

import (
	"auth-service/models"
	"auth-service/utils"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	var req models.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	_, err = db.Exec("INSERT INTO users (email, username, hashed_password) VALUES ($1, $2, $3)", req.Email, req.Username, hashedPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User registered successfully"})
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var storedPassword string
	var userID string
	err := db.QueryRow("SELECT id, hashed_password FROM users WHERE email=$1", req.Email).Scan(&userID, &storedPassword)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
	}

	return c.Status(fiber.StatusOK).JSON(models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// RefreshToken handles the refreshing of the access token
func RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing refresh token request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	log.Println("Received refresh token request with token:", req.RefreshToken)

	claims, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		log.Println("Error parsing refresh token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		log.Println("Invalid token claims, no 'sub' field found:", claims)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	log.Println("Generating new access token for user ID:", userID)

	// Generate a new access token
	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		log.Println("Error generating new access token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token generation failed"})
	}

	log.Println("Successfully generated new access token for user ID:", userID)

	return c.Status(fiber.StatusOK).JSON(models.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: req.RefreshToken,
	})
}

// Logout invalidates the user's session (client should clear the tokens)
func Logout(c *fiber.Ctx) error {
	// Simply return a success response (client clears tokens)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}
