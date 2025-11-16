package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PostgresDB wraps the database connection
type PostgresDB struct {
	DB *sqlx.DB
}

// NewPostgresDB creates a new PostgreSQL connection
func NewPostgresDB(dsn string) (*PostgresDB, error) {
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("✓ PostgreSQL connected successfully")

	return &PostgresDB{DB: db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	return p.DB.Close()
}
