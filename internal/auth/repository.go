package auth

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	SaveRefreshToken(userID, refreshToken string) error
	DeleteRefreshToken(userID string) error
	CheckRefreshToken(userID, refreshToken string) error
	FindByToken(token string) (string, error)
}

type PostgresUserRepo struct {
	ctx context.Context
	db  *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &PostgresUserRepo{
		ctx: context.Background(),
		db:  db,
	}
}

func (r *PostgresUserRepo) FindByEmail(email string) (*User, error) {
	row := r.db.QueryRow("SELECT id, email, password FROM users WHERE email=$1", email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not find user by email: %w", err)
	}
	return user, nil
}

func (r *PostgresUserRepo) Create(user *User) error {
	return r.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email, user.Password,
	).Scan(&user.ID)
}

func (r *PostgresUserRepo) SaveRefreshToken(userID, refreshToken string) error {

	query := `insert into refresh_tokens (user_id, token, expires_at) values ($1, $2, NOW() + INTERVAL '7 days')`
	_, err := r.db.Exec(query, userID, "Bearer "+refreshToken)
	if err != nil {
		return fmt.Errorf("could not save refresh token: %w", err)
	}
	return nil
}

func (r *PostgresUserRepo) DeleteRefreshToken(userID string) error {

	query := `delete from refresh_tokens where user_id=$1`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("could not delete refresh token: %w", err)
	}
	return nil
}

func (r *PostgresUserRepo) CheckRefreshToken(userID, refreshToken string) error {

	var exists bool
	query := `select exists(select 1 from refresh_tokens where user_id=$1 and token=$2 and expires_at > NOW())`
	err := r.db.QueryRow(query, userID, refreshToken).Scan(&exists)
	if err != nil {
		return fmt.Errorf("could not check refresh token: %w", err)
	}
	if !exists {
		return fmt.Errorf("refresh token does not exist")
	}
	return nil
}

func (r *PostgresUserRepo) FindByToken(token string) (string, error) {

	query := `select user_id from refresh_tokens where token=$1`
	row := r.db.QueryRow(query, token)
	var userID string
	err := row.Scan(&userID)
	if err != nil {
		return "", fmt.Errorf("could not find user by token: %w", err)
	}
	return userID, nil
}
