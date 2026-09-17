package postgres

import (
	"database/sql"
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// RunMigrations runs database migrations located at migrationsPath (e.g. "file://path/to/migrations").
func RunMigrations(databaseUrl, migrationsPath string) error {
	if databaseUrl == "" {
		return fmt.Errorf("database url vazio")
	}

	db, err := sql.Open("pgx", databaseUrl)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	return nil
}
