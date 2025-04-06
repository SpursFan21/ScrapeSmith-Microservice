//user-service\models\user.go
package models

import "database/sql"

type User struct {
	ID             string         `json:"id"`
	Username       string         `json:"username"`
	Email          string         `json:"email"`
	Name           sql.NullString `json:"name"`
	Image          sql.NullString `json:"image"`
	HashedPassword string         `json:"-"`
	CreatedAt      string         `json:"created_at"`
}
