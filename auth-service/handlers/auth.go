// auth-service\handlers\auth.go
package handlers

import (
	"auth-service/database"
	"auth-service/models"
	"auth-service/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *fiber.Ctx) error {
	var req models.SignupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	//  Check if email or username already exists
	existsFilter := bson.M{"$or": []bson.M{
		{"email": req.Email},
		{"username": req.Username},
	}}

	count, err := database.UserCollection.CountDocuments(context.TODO(), existsFilter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email or username already in use"})
	}

	// Continue with creation
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Hash error"})
	}

	user := models.User{
		Email:          req.Email,
		Username:       req.Username,
		HashedPassword: string(hashed),
		IsAdmin:        false,
		CreatedAt:      time.Now(),
	}

	_, err = database.UserCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User registered"})
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var user models.User
	err := database.UserCollection.FindOne(context.TODO(), bson.M{"email": req.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	accessToken, _ := utils.GenerateAccessToken(user.ID.Hex(), user.IsAdmin)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID.Hex(), user.IsAdmin)

	return c.JSON(models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		IsAdmin:      user.IsAdmin,
	})
}

func RefreshToken(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err != nil || req.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid refresh request"})
	}

	claims, err := utils.ParseToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	isAdmin, _ := claims["is_admin"].(bool)

	accessToken, err := utils.GenerateAccessToken(userID, isAdmin)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token generation failed"})
	}

	return c.JSON(models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: req.RefreshToken, // keep same refresh token
		IsAdmin:      isAdmin,
	})
}

func Logout(c *fiber.Ctx) error {
	// For stateless JWT, logout is handled client-side (clearing tokens)
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}
