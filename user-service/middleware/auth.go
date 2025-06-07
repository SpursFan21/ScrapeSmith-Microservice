// user-service\middleware\auth.go

package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		TokenLookup:  "header:Authorization",
		AuthScheme:   "Bearer",
		ContextKey:   "jwt",
		Claims:       &jwt.MapClaims{},
		ErrorHandler: jwtErrorHandler,
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("jwt").(*jwt.Token)
			claims := user.Claims.(*jwt.MapClaims)

			sub, ok := (*claims)["sub"].(string)
			if !ok {
				log.Println("JWT SuccessHandler: invalid sub claim in token")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid token structure",
				})
			}

			log.Printf("JWT SuccessHandler: extracted userId %s", sub)
			c.Locals("userId", sub)
			return c.Next()
		},
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	log.Printf("JWT ErrorHandler triggered: %v", err)

	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
}
