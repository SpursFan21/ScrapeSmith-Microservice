//ScrapeSmith\payment-service\utils\mongo.go

package utils

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDatabase *mongo.Database

func InitMongo() {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping failed: %v", err)
	}

	log.Println("✅ Connected to MongoDB")

	MongoClient = client
	MongoDatabase = client.Database(dbName)
}

func GetCollection(name string) *mongo.Collection {
	return MongoDatabase.Collection(name)
}
