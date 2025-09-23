// middleware/jwt.go
package middleware

import (
	"strings"

	"github.com/OzyKleyton/studio-api/utils/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthJwt() fiber.Handler {
	return func(c *fiber.Ctx) error {

		if c.Method() == "OPTIONS" {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Token de autenticação é necessário",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := auth.VerifyToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Unauthorized",
				"message": "Token inválido ou expirado: " + err.Error(),
			})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("claims", claims)

		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) uint {
	if userID, ok := c.Locals("userID").(uint); ok {
		return userID
	}
	return 0
}

func GetUsername(c *fiber.Ctx) string {
	if username, ok := c.Locals("username").(string); ok {
		return username
	}
	return ""
}

func GetEmail(c *fiber.Ctx) string {
	if email, ok := c.Locals("email").(string); ok {
		return email
	}
	return ""
}
