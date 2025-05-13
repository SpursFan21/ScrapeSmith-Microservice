// ScrapeSmith\auth-service\database\mongo.go
package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var UserCollection *mongo.Collection

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(os.Getenv("MONGO_URI"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Mongo connection error: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Mongo ping failed: %v", err)
	}

	Client = client
	dbName := os.Getenv("MONGO_DB")
	UserCollection = client.Database(dbName).Collection("users")

	//  Create unique indexes on email and username
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	if _, err := UserCollection.Indexes().CreateMany(ctx, indexModels); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}

	fmt.Println("✅ auth-service connected to MongoDB")
}
