package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

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
	query := `
		SELECT id, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1 AND is_active = true
	`

	user := &domain.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

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

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil
}

// EmailExists checks if an email is already registered
func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}
