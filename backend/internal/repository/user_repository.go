package repository

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete soft deletes a user (sets IsActive to false)
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users with pagination
	List(ctx context.Context, offset, limit int) ([]*domain.User, error)

	// EmailExists checks if an email is already registered
	EmailExists(ctx context.Context, email string) (bool, error)

	// GetByEmailAll retrieves a user by email regardless of status
	GetByEmailAll(ctx context.Context, email string) (*domain.User, error)
}
