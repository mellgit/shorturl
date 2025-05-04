package shortener

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	Save(url *URL) error
	IsAliasTaken(alias string) (bool, error)
	Stats(alias string) (int, error)
	List() (*[]URL, error)
	Delete(alias string) error
}

type PostgresRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) Repository {
	return &PostgresRepo{db}
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
		return false, err
	}
	return exists, nil
}

func (r *PostgresRepo) Stats(alias string) (int, error) {
	query := `select count(*) from clicks where alias=$1;`
	rows, err := r.db.Query(query, alias)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
	}
	return count, nil

}

func (r *PostgresRepo) List() (*[]URL, error) {
	ctx := context.Background()
	query := `select * from urls;`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var urls []URL
	for rows.Next() {
		var u URL
		if err := rows.Scan(&u.ID, &u.UserID, &u.Original, &u.Alias, &u.ExpiresAt, &u.CreatedAt); err != nil {
			return nil, err
		}
		urls = append(urls, u)

	}
	return &urls, nil

}

func (r *PostgresRepo) Delete(alias string) error {

	ctx := context.Background()
	query := `delete from urls where alias=$1;`
	_, err := r.db.ExecContext(ctx, query, alias)
	if err != nil {
		return err
	}
	return nil
}
