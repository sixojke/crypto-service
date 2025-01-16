package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/sixojke/crypto-service/internal/config"
	log "github.com/sixojke/crypto-service/pkg/logger"
)

const (
	maxAttempts = 5
	retryDelay  = time.Second
)

// NewPostgresDB creates a new connection to the Postgres database.
func NewPostgresDB(cfg config.Postgres) (*sqlx.DB, error) {
	dsn := buildPostgresDSN(cfg)

	var db *sqlx.DB
	var err error

	for i := 0; i < maxAttempts; i++ {
		db, err = sqlx.Open("postgres", dsn)
		if err == nil {
			break
		}

		if i < maxAttempts-1 {
			log.Warnf("failed to connect to db (attempt %d), retrying in %v: %v", i+1, retryDelay, err)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to db after multiple retries: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err := migratePostgres(db, cfg.MigrationsPath); err != nil {
		return nil, fmt.Errorf("migration error: %w", err)
	}

	return db, nil
}

// migratePostgres applies migrations to the Postgres.
func migratePostgres(db *sqlx.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver)

	if err != nil {
		return fmt.Errorf("creating migration instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	log.Info("migrations applied successfully")
	return nil
}

// buildPostgresDSN creates a PostgreSQL connection method
func buildPostgresDSN(cfg config.Postgres) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
}
