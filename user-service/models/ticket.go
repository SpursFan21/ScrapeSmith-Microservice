//ScrapeSmith\user-service\models\ticket.go

package models

import "time"

type TicketResponse struct {
	FromAdmin bool      `json:"fromAdmin,omitempty" bson:"fromAdmin,omitempty"`
	AdminID   string    `json:"adminId,omitempty"   bson:"adminId,omitempty"`
	Message   string    `json:"message"             bson:"message"`
	Timestamp time.Time `json:"timestamp"           bson:"timestamp"`
}

type Ticket struct {
	ID        string           `json:"id" bson:"_id,omitempty"`
	UserID    string           `json:"userId" bson:"userId"`
	Subject   string           `json:"subject" bson:"subject"`
	Message   string           `json:"message" bson:"message"`
	Status    string           `json:"status" bson:"status"` // open or closed
	Responses []TicketResponse `json:"responses" bson:"responses"`
	CreatedAt time.Time        `json:"createdAt" bson:"createdAt"`
}
