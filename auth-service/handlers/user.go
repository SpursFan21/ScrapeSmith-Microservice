// ScrapeSmith\auth-service\handlers\user.go
package handlers

import (
	"auth-service/database"
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("user").(string)
	oid, _ := primitive.ObjectIDFromHex(userID)

	var user struct {
		Username string `bson:"username"`
		Email    string `bson:"email"`
	}
	err := database.UserCollection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&user)
	if err != nil {
		log.Println("User lookup error:", err)
		return c.Status(500).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"id":       userID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func UpdateUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("user").(string)
	oid, _ := primitive.ObjectIDFromHex(userID)

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid body"})
	}

	_, err := database.UserCollection.UpdateOne(context.TODO(),
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"username": req.Username, "email": req.Email}},
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update"})
	}

	return c.JSON(fiber.Map{"message": "Profile updated"})
}
