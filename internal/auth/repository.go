package auth

import "database/sql"

type UserRepository interface {
	FindByEmail(email string) (*User, error)
	Create(user *User) error
}

type PostgresUserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db}
}

func (r *PostgresUserRepo) FindByEmail(email string) (*User, error) {
	row := r.db.QueryRow("SELECT id, email, password FROM users WHERE email=$1", email)
	user := &User{}
	err := row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepo) Create(user *User) error {
	return r.db.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id",
		user.Email, user.Password,
	).Scan(&user.ID)
}
