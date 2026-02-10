package middlewares

import (
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticateJWT(c *fiber.Ctx) error {
	cookie := c.Get("Authorization") // Biasanya "Bearer <token>"
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Unauthenticated"})
	}

	tokenString := cookie[7:] // Menghapus tulisan "Bearer "
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Token tidak valid"})
	}

	return c.Next()
}