package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
)

// RunMigrations executes database migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Get absolute path to migrations directory
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to resolve migrations path: %w", err)
	}

	// Convert Windows backslashes to forward slashes
	absPath = filepath.ToSlash(absPath)

	// Create postgres driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	// Format: file:// + absolute path (with forward slashes)
	// For Windows: file://D:/Repos/Arnela/backend/migrations
	// For Unix: file:///path/to/migrations
	sourceURL := fmt.Sprintf("file://%s", absPath)

	log.Info().Str("sourceURL", sourceURL).Msg("Attempting to run migrations")

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Info().Msg("Migrations applied successfully")
	return nil
}

// pathToFileURL converts a filesystem path to a proper file:// URL for golang-migrate
func pathToFileURL(absPath string) string {
	// Normalize path separators to forward slashes
	normalizedPath := filepath.ToSlash(absPath)

	// For Windows absolute paths (C:/path), use file:///
	// For Unix absolute paths (/path), use file://
	if len(normalizedPath) > 1 && normalizedPath[1] == ':' {
		// Windows: file:///C:/path/to/migrations
		return "file:///" + normalizedPath
	}

	// Unix: file:///path/to/migrations
	return "file://" + normalizedPath
}
