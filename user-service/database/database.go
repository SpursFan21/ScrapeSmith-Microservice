//user-service\database\database.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect to the PostgreSQL database
func Connect() {
	// Load database connection string from environment variables
	dbURI := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)

	// Open the database connection
	var err error
	DB, err = sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping the database to verify the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	// Optional: Set connection pool limits (for better performance)
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * 60 * 1000)

	// Log successful connection
	fmt.Println("✅ user-service Connected to PostgreSQL")
}
