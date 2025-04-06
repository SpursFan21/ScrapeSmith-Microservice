//user-service\database\mongo.go
package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

func ConnectMongo() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI not set in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	MongoClient = client
	log.Println("✅ user-service connected to MongoDB")
}

func GetMongoCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("MONGO_DB")
	return MongoClient.Database(dbName).Collection(collectionName)
}
