package redirect

import (
	"database/sql"
	"time"
)

type Repository interface {
	FindOriginalByAlias(alias string) (string, time.Time, error)
	SaveClick(c *Click) error
}
type PostgresRedirectRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &PostgresRedirectRepo{db}
}

func (r *PostgresRedirectRepo) FindOriginalByAlias(alias string) (string, time.Time, error) {
	var original string
	var expiresAt time.Time
	err := r.db.QueryRow(
		"SELECT original, expires_at FROM urls WHERE alias = $1", alias,
	).Scan(&original, &expiresAt)
	return original, expiresAt, err
}

func (r *PostgresRedirectRepo) SaveClick(c *Click) error {
	_, err := r.db.Exec(
		"INSERT INTO clicks (alias, ip, user_agent, created_at) VALUES ($1, $2, $3, $4)",
		c.Alias, c.IP, c.UserAgent, time.Now(),
	)
	return err
}
