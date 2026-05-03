package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/w0ikid/highload-auth-go/pkg/crypto/jwt"
)

func AuthMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
		}

		tokenString := parts[1]
		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid or expired token"})
		}

		// Store user_id in context
		c.Locals("user_id", claims.UserID)
		return c.Next()
	}
}
