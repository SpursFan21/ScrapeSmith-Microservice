// user-service\handlers\users.go
// user-service\handlers\users.go
package handlers

import (
	"context"
	"log"
	"time"
	"user-service/models"
	"user-service/mongo"
	"user-service/utils"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetUser fetches the user's account data
func GetUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = mongo.GetCollection("users").FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"id":         user.ID.Hex(),
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name,
		"image":      user.Image,
		"created_at": user.CreatedAt,
	})
}

// UpdateUser updates the user's account data
func UpdateUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Image    string `json:"image"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		log.Println("Error parsing update request:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	update := bson.M{"$set": bson.M{
		"username": updateData.Username,
		"email":    updateData.Email,
		"name":     updateData.Name,
		"image":    updateData.Image,
	}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = mongo.GetCollection("users").UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		log.Println("Error updating user data:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}

// UpdatePassword securely updates the user's password
func UpdatePassword(c *fiber.Ctx) error {
	userID := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = mongo.GetCollection("users").FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if !utils.CheckPassword(user.HashedPassword, input.OldPassword) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect current password"})
	}

	newHashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		log.Println("Error hashing new password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	_, err = mongo.GetCollection("users").UpdateOne(ctx,
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"hashed_password": newHashedPassword}},
	)
	if err != nil {
		log.Println("Error updating password:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}
