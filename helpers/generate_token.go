package helpers

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id uint, expirationTime time.Time, c *fiber.Ctx) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": id,
		"exp": expirationTime.Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error creating JWT token"})
	}
	return tokenString, nil
}
