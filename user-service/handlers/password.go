package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"user-service/utils"

	"github.com/gofiber/fiber/v2"
)

// UpdatePassword updates the user's password securely
func UpdatePassword(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Params("id")

	var input struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		log.Println("Failed to parse request body:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var storedHash string
	err := db.QueryRow(`SELECT hashed_password FROM users WHERE id = $1`, userID).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found when attempting password update:", userID)
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		log.Println("Error fetching hashed password:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	// Verify the old password
	if !utils.CheckPassword(storedHash, input.OldPassword) {
		log.Println("Incorrect current password for user ID:", userID)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Incorrect current password"})
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		log.Println("Error hashing new password:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error hashing new password"})
	}

	// Update the hashed password in the database
	_, err = db.Exec(`UPDATE users SET hashed_password = $1 WHERE id = $2`, hashedPassword, userID)
	if err != nil {
		log.Println("Error updating password for user ID:", userID, "Error:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update password"})
	}

	log.Println("Password updated successfully for user ID:", userID)
	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}
