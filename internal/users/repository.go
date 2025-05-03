package users

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	ListUsers() (*[]UserResponse, error)
}

type PostgresRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &PostgresRepo{db}
}

func (r *PostgresRepo) ListUsers() (*[]UserResponse, error) {

	ctx := context.Background()
	query := `SELECT * FROM users`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []UserResponse
	for rows.Next() {

		var user UserResponse
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to get all users: %w", err)
		}
		users = append(users, user)
	}
	return &users, nil
}
