package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
		}

		token, err := parseToken(tokenStr, false)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"].(string))

		return c.Next()
	}
}

func parseToken(tokenStr string, isRefresh bool) (*jwt.Token, error) {
	secret := os.Getenv("ACCESS_KEY") // secret token
	if isRefresh {
		secret = os.Getenv("REFRESH_KEY")
	}

	// delete Bearer prefix before transferring the token to jwt.Parse
	parts := strings.Split(tokenStr, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("invalid token")
	}

	// get jwt without prefix
	tokenStr = parts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
