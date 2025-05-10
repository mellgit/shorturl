package middleware

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
	"time"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
		}

		token, err := ParseToken(tokenStr, false)
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
		}

		claims := token.Claims.(jwt.MapClaims)
		c.Locals("user_id", claims["user_id"].(string))

		return c.Next()
	}
}

func ParseToken(tokenStr string, isRefresh bool) (*jwt.Token, error) {
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

// GenerateToken - generate token. isRefresh false - access, isRefresh true - refresh
func GenerateToken(userID string, isRefresh bool) (string, error) {

	expirationTime := time.Now().Add(5 * time.Minute) // access token on 5 min
	secretKey := os.Getenv("ACCESS_KEY")
	if isRefresh {
		expirationTime = time.Now().Add(7 * 24 * time.Hour) // refresh token on 1 week
		secretKey = os.Getenv("REFRESH_KEY")
	}
	// data for token
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // create new token (algorithm signing HMAC-SHA256)

	signToken, err := token.SignedString([]byte(secretKey)) // header.payload.signature
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return signToken, nil
}
