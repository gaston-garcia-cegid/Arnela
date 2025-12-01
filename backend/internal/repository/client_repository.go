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

	// GetByDNICIF retrieves a client by DNI/CIF (excluding soft-deleted)
	GetByDNICIF(ctx context.Context, dniCif string) (*domain.Client, error)

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

	// DNICIFExists checks if a DNI/CIF is already registered (excluding soft-deleted)
	DNICIFExists(ctx context.Context, dniCif string, excludeID *uuid.UUID) (bool, error)
}

// ClientFilters represents filters for listing clients
type ClientFilters struct {
	Search   string // Search in name, email, phone, DNI/CIF
	IsActive *bool  // Filter by active status
	City     string // Filter by city
	Province string // Filter by province
}
