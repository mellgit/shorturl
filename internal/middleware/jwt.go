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

func GenerateAccessToken(userID string) (string, error) {

	expirationTime := time.Now().Add(5 * time.Minute) // access token on 5 min
	// generate access token
	// data token
	accClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accClaims) // create new token (algorithm signing HMAC-SHA256)
	accSecret := os.Getenv("ACCESS_KEY")                             // secret token

	// use secret key for sing token
	accessToken, err := accToken.SignedString([]byte(accSecret)) // header.payload.signature
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return accessToken, nil
}

func GenerateRefreshToken(userID string) (string, error) {

	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)
	// generate refresh token
	refClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     refreshExpiration.Unix(),
	}

	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)
	refSecret := os.Getenv("REFRESH_KEY")

	refreshToken, err := refToken.SignedString([]byte(refSecret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}
	return refreshToken, nil
}
