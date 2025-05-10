package auth

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
	SaveRefreshToken(user *User, refreshToken string) error
	DeleteRefreshToken(user *User) error
	CheckRefreshToken(userID, refreshToken string) error
	FindByToken(token string) (*User, error)
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

func (r *PostgresUserRepo) SaveRefreshToken(user *User, refreshToken string) error {

	query := `update users set token=$1, expires_at=NOW() + INTERVAL '7 days' where id=$2`
	_, err := r.db.Exec(query, "Bearer "+refreshToken, user.ID)
	if err != nil {
		return fmt.Errorf("could not save refresh token: %w", err)
	}
	return nil
}

func (r *PostgresUserRepo) DeleteRefreshToken(user *User) error {

	query := `update users set token=null where id=$1;`
	_, err := r.db.Exec(query, user.ID)
	if err != nil {
		return fmt.Errorf("could not delete refresh token: %w", err)
	}
	return nil
}

func (r *PostgresUserRepo) CheckRefreshToken(userID string, refreshToken string) error {

	var exists bool
	query := `select exists(select 1 from users where id=$1 and token=$2 and expires_at > NOW())`
	err := r.db.QueryRow(query, userID, refreshToken).Scan(&exists)
	if err != nil {
		return fmt.Errorf("could not check refresh token: %w", err)
	}
	if !exists {
		return fmt.Errorf("refresh token does not exist")
	}
	return nil
}

func (r *PostgresUserRepo) FindByToken(token string) (*User, error) {

	query := `SELECT id, email, password FROM users WHERE token=$1`
	row := r.db.QueryRow(query, token)
	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not find user by token: %w", err)
	}
	return user, nil
}
