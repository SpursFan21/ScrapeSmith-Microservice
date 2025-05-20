//ScrapeSmith\payment-service\models\user_balance.go
// forge balance

package models

import "time"

type UserBalance struct {
	UserID      string    `bson:"user_id" json:"user_id"`
	Balance     int       `bson:"balance" json:"balance"`
	LastUpdated time.Time `bson:"last_updated" json:"last_updated"`
}
