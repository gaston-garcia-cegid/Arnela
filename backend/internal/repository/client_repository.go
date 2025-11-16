package repository

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// ClientRepository defines the interface for client data operations
type ClientRepository interface {
	// Create creates a new client
	Create(ctx context.Context, client *domain.Client) error

	// GetByID retrieves a client by ID (excluding soft-deleted)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error)

	// GetByEmail retrieves a client by email (excluding soft-deleted)
	GetByEmail(ctx context.Context, email string) (*domain.Client, error)

	// GetByDNI retrieves a client by DNI (excluding soft-deleted)
	GetByDNI(ctx context.Context, dni string) (*domain.Client, error)

	// GetByUserID retrieves a client by associated user ID (excluding soft-deleted)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error)

	// Update updates a client's information
	Update(ctx context.Context, client *domain.Client) error

	// Delete soft-deletes a client
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a paginated list of clients with optional filters
	List(ctx context.Context, filters ClientFilters, offset, limit int) ([]*domain.Client, error)

	// Count returns the total number of clients matching the filters
	Count(ctx context.Context, filters ClientFilters) (int, error)

	// EmailExists checks if an email is already registered (excluding soft-deleted)
	EmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)

	// NIFExists checks if a NIF is already registered (excluding soft-deleted)
	NIFExists(ctx context.Context, nif string, excludeID *uuid.UUID) (bool, error)

	// DNIExists checks if a DNI is already registered (excluding soft-deleted)
	DNIExists(ctx context.Context, dni string, excludeID *uuid.UUID) (bool, error)

	// GetByNIF retrieves a client by NIF (excluding soft-deleted)
	GetByNIF(ctx context.Context, nif string) (*domain.Client, error)
}

// ClientFilters represents filters for listing clients
type ClientFilters struct {
	Search   string // Search in name, email, phone, DNI
	IsActive *bool  // Filter by active status
	City     string // Filter by city
	Province string // Filter by province
}
