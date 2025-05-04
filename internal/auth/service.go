package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(email, password string) error
	Login(email, password string) (string, error)
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

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", fmt.Errorf("could not find user by email: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("could not compare password: %w", err)
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("could not sign token: %w", err)
	}

	return signedToken, nil
}
