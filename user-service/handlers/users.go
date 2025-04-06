//user-service\handlers\users.go
package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"user-service/models"

	"github.com/gofiber/fiber/v2"
)

// GetUser fetches the user's account data
func GetUser(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Params("id")

	var user models.User
	err := db.QueryRow(`
        SELECT id, username, email, name, image, created_at
        FROM users
        WHERE id = $1`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.Image, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		log.Printf("Error fetching user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	response := fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name.String,
		"image":      user.Image.String,
		"created_at": user.CreatedAt,
	}

	return c.JSON(response)
}

// UpdateUser updates the user's account data
func UpdateUser(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Params("id")

	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Image    string `json:"image"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		log.Println("Error parsing update request:", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	_, err := db.Exec(`
		UPDATE users
		SET username = $1, email = $2, name = $3, image = $4
		WHERE id = $5`,
		updateData.Username, updateData.Email, updateData.Name, updateData.Image, userID)

	if err != nil {
		log.Println("Error updating user data:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}
