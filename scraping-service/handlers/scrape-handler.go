// scraping-service\handlers\scrape-handler.go
package handlers

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"scraping-service/models"
	"scraping-service/utils"
)

type ScrapeRequest struct {
	URL        string `json:"url"`
	Analysis   string `json:"analysis_type"`
	CustomCode string `json:"custom_script,omitempty"`
}

func SingleScrape(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("Missing Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid Authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Println("JWT_SECRET_KEY missing from env")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server misconfiguration: missing JWT_SECRET_KEY",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("JWT parse error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to parse token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	userIDVal, ok := claims["sub"]
	if !ok {
		log.Println("user_id (sub) missing in token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user_id missing in token claims",
		})
	}

	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		log.Println("Invalid user_id in token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in token claims",
		})
	}

	var req ScrapeRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	orderId := uuid.New().String()
	createdAt := time.Now()

	job := models.QueuedScrapeJob{
		OrderID:      orderId,
		UserID:       userID,
		CreatedAt:    createdAt,
		URL:          req.URL,
		AnalysisType: req.Analysis,
		CustomScript: req.CustomCode,
		Status:       "pending",
		Attempts:     0,
	}

	collection := utils.GetCollection("queued_scrape_jobs")
	_, err = collection.InsertOne(c.Context(), job)
	if err != nil {
		log.Printf("Failed to enqueue scrape job: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to enqueue scrape job",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":    "Job received and queued for scraping",
		"order_id":   orderId,
		"created_at": createdAt,
	})
}
