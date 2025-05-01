package auth

import (
	"errors"
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
		return errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
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
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
