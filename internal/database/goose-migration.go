package database

import (
	"book-catalog/internal/config"
	"book-catalog/internal/logger"
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" //nolint:revive // required for goose
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// GooseMigrateDatabase migrates the database using goose
func GooseMigrateDatabase(cfg *config.Config, log logger.Logger) error {
	log.Trc().Msg("parse config for database migration")

	// init db
	db, err := sql.Open("pgx/v5", cfg.ConnDB)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// check connection
	log.Trc().Msg("ping database")
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// setup
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(log)
	goose.SetVerbose(true)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// run
	log.Trc().Msg("run migrations")
	return goose.Up(db, "migrations")
}
