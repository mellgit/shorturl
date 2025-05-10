package redirect

import (
	"database/sql"
	"fmt"
	"time"
)

type PostgresRepository interface {
	FindOriginalByAlias(alias string) (string, time.Time, error)
	SaveClick(c *Click) error
}
type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) PostgresRepository {
	return &postgresRepository{db}
}

func (r *postgresRepository) FindOriginalByAlias(alias string) (string, time.Time, error) {
	var original string
	var expiresAt time.Time
	err := r.db.QueryRow(
		"SELECT original, expires_at FROM urls WHERE alias = $1", alias,
	).Scan(&original, &expiresAt)
	return original, expiresAt, fmt.Errorf("could not find original alias: %v", err)
}

func (r *postgresRepository) SaveClick(c *Click) error {
	_, err := r.db.Exec(
		"INSERT INTO clicks (alias, ip, user_agent, created_at) VALUES ($1, $2, $3, $4)",
		c.Alias, c.IP, c.UserAgent, time.Now(),
	)
	return fmt.Errorf("could not save click analytic: %v", err)
}
