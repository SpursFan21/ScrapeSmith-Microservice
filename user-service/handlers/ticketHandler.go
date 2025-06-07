//ScrapeSmith\user-service\handlers\ticketHandler.go

package handlers

import (
	"context"
	"log"
	"time"

	"user-service/models"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ticketColl *mongo.Collection

func InitTicketCollection(c *mongo.Client, dbName string) {
	ticketColl = c.Database(dbName).Collection("tickets")
}

// POST /users/tickets - submit a ticket
func SubmitTicket(c *fiber.Ctx) error {
	userIDVal := c.Locals("userId")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		log.Println("Invalid or missing userId in context")
		return c.Status(400).JSON(fiber.Map{"error": "Unauthorized or missing user ID"})
	}

	var payload struct {
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	ticket := models.Ticket{
		UserID:    userID,
		Subject:   payload.Subject,
		Message:   payload.Message,
		Status:    "open",
		Responses: []models.TicketResponse{},
		CreatedAt: time.Now(),
	}

	_, err := ticketColl.InsertOne(context.TODO(), ticket)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit ticket"})
	}

	log.Printf("Ticket submitted by user %s | Subject: %s", userID, payload.Subject)

	return c.JSON(fiber.Map{"message": "Ticket submitted"})
}

// GET /users/tickets - get user's tickets
func GetMyTickets(c *fiber.Ctx) error {
	userIDVal := c.Locals("userId")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		log.Println(" Invalid or missing userId in GetMyTickets")
		return c.Status(400).JSON(fiber.Map{"error": "Unauthorized or missing user ID"})
	}

	//  Query with string
	cursor, err := ticketColl.Find(context.TODO(), bson.M{"userId": userID})
	if err != nil {
		log.Printf(" Failed to fetch tickets for user %s: %v", userID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tickets"})
	}

	var tickets []models.Ticket
	if err := cursor.All(context.TODO(), &tickets); err != nil {
		log.Printf(" Failed to parse tickets for user %s: %v", userID, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to parse tickets"})
	}

	log.Printf(" Fetched %d tickets for user %s", len(tickets), userID)
	return c.JSON(tickets)
}

// POST /users/tickets/:id/reply - reply to a ticket
func ReplyToTicket(c *fiber.Ctx) error {
	userIDVal := c.Locals("userId")
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		log.Println("Invalid or missing userId in ReplyToTicket")
		return c.Status(400).JSON(fiber.Map{"error": "Unauthorized or missing user ID"})
	}

	ticketID := c.Params("id")
	objID, err := primitive.ObjectIDFromHex(ticketID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ticket ID"})
	}

	var body struct {
		Message string `json:"message"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	update := bson.M{
		"$push": bson.M{
			"responses": models.TicketResponse{
				FromAdmin: false,
				Message:   body.Message,
				Timestamp: time.Now(),
			},
		},
	}

	result, err := ticketColl.UpdateOne(
		context.TODO(),
		bson.M{"_id": objID, "userId": userID},
		update,
	)
	if err != nil || result.MatchedCount == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "Ticket not found or unauthorized"})
	}

	log.Printf("User %s replied to ticket %s", userID, ticketID)

	return c.JSON(fiber.Map{"message": "Reply added"})
}
