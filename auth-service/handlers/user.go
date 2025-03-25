package handlers

import (
	"database/sql"
	"log"

	"github.com/gofiber/fiber/v2"
)

var db *sql.DB

// SetDB initializes the database connection (called from main)
func SetDB(database *sql.DB) {
	db = database
}

// GetUserProfile handles fetching the user profile
func GetUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("user").(string) // Get user ID from JWT claims

	var username, email string
	err := db.QueryRow("SELECT username, email FROM users WHERE id=$1", userID).Scan(&username, &email)
	if err != nil {
		log.Println("Error fetching user profile:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to fetch user profile",
		})
	}

	return c.JSON(fiber.Map{
		"id":       userID,
		"username": username,
		"email":    email,
	})
}

// UpdateUserProfile handles updating the user profile
func UpdateUserProfile(c *fiber.Ctx) error {
	type UpdateRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	userID := c.Locals("user").(string) // Get user ID from JWT claims
	var req UpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	_, err := db.Exec("UPDATE users SET username=$1, email=$2 WHERE id=$3", req.Username, req.Email, userID)
	if err != nil {
		log.Println("Error updating user profile:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to update profile",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Profile updated successfully",
	})
}
