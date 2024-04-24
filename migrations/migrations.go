package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/KretovDmitry/gophermart-loyalty-service/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

//go:embed *.sql
var fs embed.FS

func Up(db *sql.DB, cfg *config.Config) error {
	d, err := iofs.New(fs, cfg.Migrations)
	if err != nil {
		return fmt.Errorf("failed to init io/fs driver: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failde to init migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to init migrate instance: %w", err)
	}

	if err = m.Up(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}