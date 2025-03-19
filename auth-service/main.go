package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3" // Correct JWT middleware for Fiber
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

// Load environment variables
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

// Auth0 credentials from .env
var auth0Domain = os.Getenv("AUTH0_DOMAIN")
var auth0ClientID = os.Getenv("AUTH0_CLIENT_ID")
var auth0ClientSecret = os.Getenv("AUTH0_CLIENT_SECRET")

// Database connection
var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	log.Println("Database connected successfully")
}

func main() {
	initDB()
	defer db.Close()

	app := fiber.New()

	// Public route (No Auth Required)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Auth Service"})
	})

	// User Registration (Sign Up)
	app.Post("/signup", func(c *fiber.Ctx) error {
		type SignupRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Username string `json:"username"`
		}

		var req SignupRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Hash the password before saving
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to hash password"})
		}

		// Insert user into the database
		_, err = db.Exec(`
			INSERT INTO users (email, username, hashed_password) 
			VALUES ($1, $2, $3)`,
			req.Email, req.Username, hashedPassword)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating user"})
		}

		// You can also integrate Auth0 here if necessary
		// Client registration code with Auth0...

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User registered successfully"})
	})

	// User Login (Token Request)
	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Fetch the user from the database
		var storedPassword string
		var userID string
		err := db.QueryRow("SELECT id, hashed_password FROM users WHERE email=$1", req.Email).Scan(&userID, &storedPassword)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Compare the hashed password
		err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Generate JWT Token
		claims := jwt.MapClaims{
			"sub": userID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(auth0ClientSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
		}

		// Return the token
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString})
	})

	// Middleware to protect routes
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(auth0ClientSecret),
	}))

	// Protected route (Requires Auth)
	app.Get("/protected", func(c *fiber.Ctx) error {
		user := c.Locals("user") // Extract user info from token
		return c.JSON(fiber.Map{"message": "Protected Route", "user": user})
	})

	log.Println("Auth service running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
