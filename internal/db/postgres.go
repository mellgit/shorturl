package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mellgit/shorturl/internal/config"
	"github.com/pressly/goose/v3"
)

func PostgresClient(envCfg config.EnvConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		envCfg.DBHost, envCfg.DBPort, envCfg.DBName, envCfg.DBUser, envCfg.DBPassword,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}
	if err = goose.Up(db, envCfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("failed to create migrations: %w", err)
	}
	return db, nil
}
