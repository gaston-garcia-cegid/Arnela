package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies database migrations WITHOUT closing the connection
func RunMigrations(db *sql.DB, migrationsPath string) error {
	log.Printf("[DEBUG] Running migrations from: %s", migrationsPath)

	// Verificar que el directorio existe y tiene archivos
	if err := validateMigrationsDir(migrationsPath); err != nil {
		return fmt.Errorf("migrations validation failed: %w", err)
	}

	// Convertir path a URL compatible
	sourceURL := pathToMigrateURL(migrationsPath)
	log.Printf("[DEBUG] Migrations source URL: %s", sourceURL)

	// Crear driver de postgres
	driver, err := postgres.WithInstance(db, &postgres.Config{
		MigrationsTable: "schema_migrations",
		SchemaName:      "public",
	})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Crear instancia de migrate
	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// Aplicar migraciones
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		log.Println("[DEBUG] No new migrations to apply")
	} else {
		log.Println("[DEBUG] Migrations applied successfully")
	}

	// Solo cerrar el driver de migrate, NO la DB
	if serr, derr := m.Close(); serr != nil || derr != nil {
		log.Printf("[WARN] Error closing migrate instance (source=%v, db=%v)", serr, derr)
	}

	return nil
}

// validateMigrationsDir verifica que el directorio existe y tiene archivos .sql
func validateMigrationsDir(path string) error {
	// Verificar que el directorio existe
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("migrations directory does not exist: %s", path)
		}
		return fmt.Errorf("failed to stat migrations directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("migrations path is not a directory: %s", path)
	}

	// Verificar que tiene archivos .sql
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sqlFiles := 0
	log.Printf("[DEBUG] Files in migrations directory:")
	for _, entry := range entries {
		log.Printf("[DEBUG]   - %s", entry.Name())
		if !entry.IsDir() && (strings.HasSuffix(entry.Name(), ".up.sql") || strings.HasSuffix(entry.Name(), ".down.sql")) {
			sqlFiles++
		}
	}

	if sqlFiles == 0 {
		log.Printf("[WARN] No .sql migration files found in %s", path)
	} else {
		log.Printf("[DEBUG] Found %d SQL migration files", sqlFiles)
	}

	return nil
}

// pathToMigrateURL converts a filesystem path to golang-migrate compatible URL
func pathToMigrateURL(path string) string {
	// Obtener path absoluto
	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Printf("[WARN] Failed to get absolute path for %s: %v", path, err)
		absPath = path
	}

	log.Printf("[DEBUG] Absolute path: %s", absPath)
	log.Printf("[DEBUG] OS: %s", runtime.GOOS)

	// ✅ SOLUCIÓN PARA WINDOWS: golang-migrate tiene problemas con file:// en Windows
	// La solución es NO usar el prefijo file:// en Windows, solo el path
	if runtime.GOOS == "windows" {
		// Convertir backslashes a forward slashes
		absPath = filepath.ToSlash(absPath)
		log.Printf("[DEBUG] Windows path converted to: %s", absPath)

		// En Windows, golang-migrate espera: file://D:/path o simplemente D:/path
		// Pero el formato más confiable es sin file:// prefix
		return "file:" + absPath
	}

	// En Unix/Linux/Mac, usar el formato estándar
	if !strings.HasPrefix(absPath, "/") {
		absPath = "/" + absPath
	}

	return "file://" + absPath
}
