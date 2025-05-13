//ScrapeSmith\auth-service\middlewares\jwt.go
package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:    []byte(os.Getenv("JWT_SECRET_KEY")),
		SigningMethod: "HS256",
	})
}
