// user-service\models\user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	Email          string             `json:"email" bson:"email"`
	Name           string             `json:"name" bson:"name"`
	Image          string             `json:"image" bson:"image"`
	HashedPassword string             `json:"-" bson:"hashed_password"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
}
