package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresDB creates a new PostgreSQL database connection with persistent pool
func NewPostgresDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	log.Printf("[DEBUG] Connecting to database: host=%s port=%d dbname=%s",
		cfg.Host, cfg.Port, cfg.DBName)

	// ✅ CRÍTICO: Connect establece conexión inmediatamente
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// ✅ Configuración del pool - SIN LÍMITES para evitar cierres
	db.SetMaxOpenConns(25)   // Máximo de conexiones abiertas
	db.SetMaxIdleConns(10)   // Mínimo de conexiones idle (siempre disponibles)
	db.SetConnMaxLifetime(0) // 0 = sin límite de vida
	db.SetConnMaxIdleTime(0) // 0 = nunca cerrar por inactividad

	// Verificar conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to verify database connection: %w", err)
	}

	// Log de estadísticas
	stats := db.Stats()
	log.Printf("[DEBUG] Database pool initialized - Open=%d, Idle=%d, MaxOpen=%d, MaxIdle=%d",
		stats.OpenConnections, stats.Idle, 25, 10)

	return db, nil
}

// HealthCheck verifica que el pool esté funcional y reconecta si es necesario
func HealthCheck(db *sqlx.DB) error {
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Verificar stats ANTES del ping
	stats := db.Stats()
	log.Printf("[DEBUG] HealthCheck - Pool stats BEFORE ping: Open=%d, Idle=%d, InUse=%d",
		stats.OpenConnections, stats.Idle, stats.InUse)

	// Intentar ping
	if err := db.PingContext(ctx); err != nil {
		log.Printf("[ERROR] Database ping failed: %v", err)
		return fmt.Errorf("database ping failed: %w", err)
	}

	// Verificar stats DESPUÉS del ping
	stats = db.Stats()
	log.Printf("[DEBUG] HealthCheck - Pool stats AFTER ping: Open=%d, Idle=%d, InUse=%d",
		stats.OpenConnections, stats.Idle, stats.InUse)

	if stats.OpenConnections == 0 {
		return fmt.Errorf("database pool has no open connections after ping")
	}

	return nil
}
