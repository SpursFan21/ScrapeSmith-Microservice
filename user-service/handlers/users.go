package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// GetUser fetches the user's account data
func GetUser(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Params("id") // Extract user ID from the request

	// Query the database for the user's data
	var user struct {
		ID        string         `json:"id"`
		Username  string         `json:"username"`
		Email     string         `json:"email"`
		Name      sql.NullString `json:"name"`  // Use sql.NullString for nullable fields
		Image     sql.NullString `json:"image"` // Use sql.NullString for nullable fields
		CreatedAt string         `json:"created_at"`
	}

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

	// Convert sql.NullString to string for the response
	response := fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"name":       user.Name.String,  // Convert to string
		"image":      user.Image.String, // Convert to string
		"created_at": user.CreatedAt,
	}

	return c.JSON(response)
}

// user-service\handlers\users.go
// UpdateUser updates the user's account data
func UpdateUser(c *fiber.Ctx, db *sql.DB) error {
	userID := c.Params("id") // Extract user ID from the request

	// Parse the request body
	var updateData struct {
		Username string         `json:"username"`
		Email    string         `json:"email"`
		Name     sql.NullString `json:"name"`  // Use sql.NullString for nullable fields
		Image    sql.NullString `json:"image"` // Use sql.NullString for nullable fields
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Update the user's data in the database
	_, err := db.Exec(`
		UPDATE users
		SET username = $1, email = $2, name = $3, image = $4
		WHERE id = $5`,
		updateData.Username, updateData.Email, updateData.Name, updateData.Image, userID)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
	}

	return c.JSON(fiber.Map{"message": "User updated successfully"})
}
