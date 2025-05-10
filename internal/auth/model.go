package auth

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"-"` // Храним в зашифрованном виде
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required,min=6"`
}
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
