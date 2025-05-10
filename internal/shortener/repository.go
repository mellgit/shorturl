package shortener

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	Save(url *URL) error
	IsAliasTaken(alias string) (bool, error)
	Stats(alias string) (int, error)
	List() (*[]URL, error)
	Delete(alias string) error
	UpdateAlias(alias, newAlias string) error
	GetUrlFromAlias(alias string) (string, error)
}

type PostgresRepo struct {
	db  *sql.DB
	ctx context.Context
}

func NewRepo(db *sql.DB) Repository {
	return &PostgresRepo{ctx: context.Background(), db: db}
}

func (r *PostgresRepo) Save(u *URL) error {
	query := `
	INSERT INTO urls (user_id, original, alias, expires_at, created_at)
	VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	return r.db.QueryRow(
		query,
		u.UserID, u.Original, u.Alias, u.ExpiresAt, time.Now(),
	).Scan(&u.ID)
}

func (r *PostgresRepo) IsAliasTaken(alias string) (bool, error) {
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS (SELECT 1 FROM urls WHERE alias=$1)", alias).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking if alias exists: %w", err)
	}
	return exists, nil
}

func (r *PostgresRepo) Stats(alias string) (int, error) {
	query := `select count(*) from clicks where alias=$1;`
	rows, err := r.db.Query(query, alias)
	if err != nil {
		return 0, fmt.Errorf("error getting clicks stats: %w", err)
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, fmt.Errorf("error scanning row: %w", err)
		}
	}
	return count, nil
}

func (r *PostgresRepo) List() (*[]URL, error) {
	query := `select * from urls;`
	rows, err := r.db.QueryContext(r.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listing urls: %w", err)
	}
	defer rows.Close()
	var urls []URL
	for rows.Next() {
		var u URL
		if err := rows.Scan(&u.ID, &u.UserID, &u.Original, &u.Alias, &u.ExpiresAt, &u.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		urls = append(urls, u)
	}
	return &urls, nil
}

func (r *PostgresRepo) Delete(alias string) error {
	query := `delete from urls where alias=$1;`
	_, err := r.db.ExecContext(r.ctx, query, alias)
	if err != nil {
		return fmt.Errorf("error deleting url: %w", err)
	}
	return nil
}

func (r *PostgresRepo) UpdateAlias(alias, newAlias string) error {
	query := `update urls set alias=$1 where alias=$2;`
	_, err := r.db.ExecContext(r.ctx, query, newAlias, alias)
	if err != nil {
		return fmt.Errorf("error updating alias: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetUrlFromAlias(alias string) (string, error) {
	query := `select original from urls where alias=$1;`
	row := r.db.QueryRowContext(r.ctx, query, alias)
	var original string
	if err := row.Scan(&original); err != nil {
		return "", fmt.Errorf("error scanning row: %w", err)
	}
	return original, nil
}
