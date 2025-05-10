package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	token "github.com/mellgit/shorturl/internal/middleware"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(email, password string) error
	Login(email, password string) (*TokensResponse, error)
	RefreshToken(refreshToken string) (*AccessTokenResponse, error)
	Logout(token string) error
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

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("could not find user by email: %w", err)
	}

	// compare the password hash in the database and from the user
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("could not compare password: %w", err)
	}

	accessToken, err := token.GenerateToken(user.ID.String(), false)
	if err != nil {
		return nil, fmt.Errorf("could not generate access token: %w", err)
	}
	refreshToken, err := token.GenerateToken(user.ID.String(), true)
	if err != nil {
		return nil, fmt.Errorf("could not generate refresh token: %w", err)
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

	tokenParse, err := token.ParseToken(refreshToken, true)
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	claims := tokenParse.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(string)

	// check token in db
	if err := s.repo.CheckRefreshToken(userID, refreshToken); err != nil {
		return nil, fmt.Errorf("could not check refresh token: %w", err)
	}

	accessToken, err := token.GenerateToken(userID, false)
	if err != nil {
		return nil, fmt.Errorf("could not generate access token: %w", err)
	}
	data := AccessTokenResponse{AccessToken: accessToken}
	return &data, nil
}

func (s *AuthService) Logout(tokenStr string) error {

	user, err := s.repo.FindByToken(tokenStr)
	if err != nil {
		return fmt.Errorf("could not find user by token: %w", err)
	}
	if err := s.repo.DeleteRefreshToken(user); err != nil {
		return fmt.Errorf("could not delete refresh token: %w", err)
	}
	return nil
}
