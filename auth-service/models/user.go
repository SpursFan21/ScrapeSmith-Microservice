// ScrapeSmith\auth-service\models\user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Email          string             `bson:"email"`
	Username       string             `bson:"username"`
	HashedPassword string             `bson:"hashed_password"`
	IsAdmin        bool               `bson:"is_admin"`
	CreatedAt      time.Time          `bson:"created_at"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
