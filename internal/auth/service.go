package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(email, password string) error
	Login(email, password string) (*TokensResponse, error)
	RefreshToken(refreshToken string) (*AccessTokenResponse, error)
}
type AuthService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &AuthService{repo}
}

func (s *AuthService) Register(email, password string) error {
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return fmt.Errorf("email already registered")
	}

	// get hash from password for save in db
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("could not hash password: %w", err)
	}

	user := User{
		Email:    email,
		Password: string(hashed),
	}
	return s.repo.Create(&user)
}

func (s *AuthService) Login(email, password string) (*TokensResponse, error) {

	expirationTime := time.Now().Add(5 * time.Minute) // access token on 5 min
	refreshExpiration := time.Now().Add(7 * 24 * time.Hour)

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("could not find user by email: %w", err)
	}

	// compare the password hash in the database and from the user
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("could not compare password: %w", err)
	}

	// generate access token
	//data token
	accClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     expirationTime.Unix(),
	}
	accToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accClaims) // create new token (algorithm signing HMAC-SHA256)
	accSecret := os.Getenv("ACCESS_KEY")                             // secret token

	// use secret key for sing token
	accessToken, err := accToken.SignedString([]byte(accSecret)) // header.payload.signature
	if err != nil {
		return nil, fmt.Errorf("could not sign token: %w", err)
	}

	// generate refresh token
	refClaims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     refreshExpiration.Unix(),
	}

	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refClaims)
	refSecret := os.Getenv("REFRESH_KEY")

	refreshToken, err := refToken.SignedString([]byte(refSecret))
	if err != nil {
		return nil, fmt.Errorf("could not sign token: %w", err)
	}

	if err := s.repo.DeleteRefreshToken(user); err != nil {
		return nil, fmt.Errorf("could not delete refresh token: %w", err)
	}

	if err := s.repo.SaveRefreshToken(user, refreshToken); err != nil {
		return nil, fmt.Errorf("could not save refresh token: %w", err)
	}

	data := TokensResponse{AccessToken: accessToken, RefreshToken: refreshToken}

	return &data, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*AccessTokenResponse, error) {

	expirationTime := time.Now().Add(5 * time.Minute) // access token on 5 min

	tokenParse, err := s.parseToken(refreshToken, true)
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	claims := tokenParse.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	// check token in db
	if err := s.repo.CheckRefreshToken(userID, refreshToken); err != nil {
		return nil, fmt.Errorf("could not check refresh token: %w", err)
	}

	claims2 := jwt.MapClaims{
		"user_id": userID,
		"exp":     expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims2) // create new token (algorithm signing HMAC-SHA256)
	secret := os.Getenv("ACCESS_KEY")                           // secret token

	// use secret key for sing token
	accessToken, err := token.SignedString([]byte(secret)) // header.payload.signature
	if err != nil {
		return nil, fmt.Errorf("could not sign token: %w", err)
	}
	data := AccessTokenResponse{AccessToken: accessToken}
	return &data, nil
}

// todo this method duplicate in jwt package
func (s *AuthService) parseToken(tokenStr string, isRefresh bool) (*jwt.Token, error) {

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
