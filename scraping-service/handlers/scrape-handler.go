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
	// Step 1: Check for Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Println("Missing or invalid Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid Authorization header",
		})
	}

	// Step 2: Extract and validate JWT token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Println("JWT_SECRET_KEY missing from environment variables")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server misconfiguration: missing JWT_SECRET_KEY",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Printf("JWT parse error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	// Step 3: Parse claims and validate user ID
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to parse token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	userIDVal, ok := claims["sub"]
	if !ok {
		log.Println("Missing user_id (sub) in token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing user_id (sub) in token claims",
		})
	}

	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		log.Println("Invalid user_id in token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in token claims",
		})
	}

	// Step 4: Parse request body
	var req ScrapeRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if strings.TrimSpace(req.URL) == "" || strings.TrimSpace(req.Analysis) == "" {
		log.Println("Missing required fields: url or analysis_type")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: url and analysis_type",
		})
	}

	// Step 5: Create a new job
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

	// Step 6: Insert the job into the queue
	collection := utils.GetCollection("queued_scrape_jobs")
	_, err = collection.InsertOne(c.Context(), job)
	if err != nil {
		log.Printf("Failed to enqueue scrape job: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to enqueue scrape job",
		})
	}

	// Log success
	log.Printf("✅ Scrape job queued: userId=%s, orderId=%s", userID, orderId)

	// Step 7: Return response with job details
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":   "Job received and queued for scraping",
		"orderId":   orderId,
		"createdAt": createdAt,
	})
}
