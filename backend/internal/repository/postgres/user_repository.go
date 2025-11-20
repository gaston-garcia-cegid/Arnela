package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// userRepository implements repository.UserRepository interface
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sqlx.DB) repository.UserRepository {
	if db == nil {
		panic("database connection is nil")
	}
	log.Printf("[DEBUG] UserRepository initialized with db pointer: %p", db)
	return &userRepository{db: db}
}

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	user := &domain.User{}
	// Use GetContext instead of QueryRowContext for better error handling
	err := r.db.GetContext(ctx, user, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Log de debug
	log.Printf("[DEBUG] GetByEmail called for: %s", email)
	log.Printf("[DEBUG] Repository db pointer: %p", r.db)

	// Verificar que el puntero no sea nil
	if r.db == nil {
		log.Printf("[ERROR] Database connection is nil in repository!")
		return nil, fmt.Errorf("database connection is nil")
	}

	// Verificar estado del pool
	stats := r.db.Stats()
	log.Printf("[DEBUG] DB Pool stats before query - Open=%d, Idle=%d, InUse=%d, WaitCount=%d",
		stats.OpenConnections, stats.Idle, stats.InUse, stats.WaitCount)

	// ✅ Verificar que el pool tiene conexiones abiertas
	if stats.OpenConnections == 0 {
		log.Printf("[ERROR] No open database connections in pool!")

		// Intentar reconectar haciendo ping
		if err := r.db.Ping(); err != nil {
			log.Printf("[ERROR] Failed to reconnect: %v", err)
			return nil, fmt.Errorf("database connection pool is empty and reconnection failed: %w", err)
		}

		log.Printf("[WARN] Successfully reconnected to database")
		stats = r.db.Stats()
		log.Printf("[DEBUG] DB Pool stats after reconnect - Open=%d, Idle=%d, InUse=%d",
			stats.OpenConnections, stats.Idle, stats.InUse)
	}

	query := `
        SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
        FROM users
        WHERE email = $1 AND is_active = true
    `

	user := &domain.User{}
	log.Printf("[DEBUG] Executing query for email: %s", email)

	// Usar GetContext en lugar de Get para respetar el contexto
	log.Printf("[DEBUG] Executing query for email: %s", email)
	err := r.db.GetContext(ctx, user, query, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("[DEBUG] User not found for email: %s", email)
			return nil, fmt.Errorf("user not found")
		}
		log.Printf("[ERROR] Query failed: %v", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	log.Printf("[DEBUG] User found: ID=%s, Email=%s, Role=%s", user.ID, user.Email, user.Role)
	return user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $2, first_name = $3, last_name = $4, role = $5, is_active = $6, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.IsActive,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete soft deletes a user by setting IsActive to false
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET is_active = false, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// List retrieves all users with pagination
func (r *userRepository) List(ctx context.Context, offset, limit int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var users []*domain.User
	// Use SelectContext for batch queries
	err := r.db.SelectContext(ctx, &users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

// EmailExists checks if an email is already registered
func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
