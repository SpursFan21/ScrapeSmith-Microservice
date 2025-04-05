// scraping-service\handlers\scrape-handler.go
package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

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
		log.Println("❌ Missing Authorization header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing or invalid Authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Println("❌ JWT_SECRET_KEY missing from env")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server misconfiguration: missing JWT_SECRET_KEY",
		})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		log.Printf("❌ JWT parse error: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("❌ Failed to parse token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	userIDVal, ok := claims["sub"]
	if !ok {
		log.Println("❌ user_id (sub) missing in token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user_id missing in token claims",
		})
	}

	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		log.Println("❌ Invalid user_id in token claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user_id in token claims",
		})
	}

	var req ScrapeRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	orderId := uuid.New().String()
	createdAt := time.Now()

	rawData, err := utils.ScrapeWithScrapeNinja(req.URL)
	if err != nil {
		log.Printf("❌ Scraping error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to scrape data",
		})
	}

	collection := utils.GetCollection("scraped_data")
	_, err = collection.InsertOne(c.Context(), map[string]interface{}{
		"order_id":      orderId,
		"user_id":       userID,
		"created_at":    createdAt,
		"url":           req.URL,
		"analysis_type": req.Analysis,
		"custom_script": req.CustomCode,
		"data":          rawData,
	})
	if err != nil {
		log.Printf("❌ Mongo insert error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save scraped data",
		})
	}

	payload := map[string]interface{}{
		"orderId":      orderId,
		"userId":       userID,
		"url":          req.URL,
		"analysisType": req.Analysis,
		"customScript": req.CustomCode,
		"createdAt":    createdAt,
		"rawData":      rawData,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("❌ Failed to marshal payload for cleaning service: %v", err)
	} else {
		go func() {
			cleanerURL := os.Getenv("DATA_CLEANER_URL")
			if cleanerURL == "" {
				cleanerURL = "http://data-cleaning-service:3004/api/clean"
			}
			log.Printf("Sending payload to cleaner: %s", cleanerURL)

			resp, err := http.Post(cleanerURL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				log.Printf("❌ Failed to send payload to cleaner: %v", err)
				return
			}
			defer resp.Body.Close()
			log.Printf("✅ Cleaner response: %s", resp.Status)
		}()
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Scraping successful",
		"order_id":   orderId,
		"created_at": createdAt,
	})
}
