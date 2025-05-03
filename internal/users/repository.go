package users

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	ListUsers() (*[]UserResponse, error)
	GetUserByID(id int64) (*UserResponse, error)
	//DeleteUserByID()
	//UpdateUserByID()
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

func (r *PostgresRepo) GetUserByID(id int64) (*UserResponse, error) {

	ctx := context.Background()
	query := `SELECT * FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	var user UserResponse
	if err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
