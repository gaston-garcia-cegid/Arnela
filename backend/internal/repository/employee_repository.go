package repository

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// EmployeeRepository defines the interface for employee data operations
type EmployeeRepository interface {
	// Create inserts a new employee
	Create(ctx context.Context, employee *domain.Employee) error

	// GetByID retrieves an employee by ID (excludes soft-deleted)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error)

	// GetByUserID retrieves an employee by their user ID
	GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error)

	// GetByEmail retrieves an employee by email
	GetByEmail(ctx context.Context, email string) (*domain.Employee, error)

	// GetByDNI retrieves an employee by DNI
	GetByDNI(ctx context.Context, dni string) (*domain.Employee, error)

	// Update updates an existing employee
	Update(ctx context.Context, employee *domain.Employee) error

	// Delete performs a soft delete
	Delete(ctx context.Context, id uuid.UUID) error

	// List returns all active employees with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.Employee, error)

	// Count returns the total number of active employees
	Count(ctx context.Context) (int, error)

	// GetBySpecialty returns employees with a specific specialty
	GetBySpecialty(ctx context.Context, specialty string) ([]*domain.Employee, error)

	// EmailExists checks if an email is already registered
	EmailExists(ctx context.Context, email string) (bool, error)

	// DNIExists checks if a DNI is already registered
	DNIExists(ctx context.Context, dni string) (bool, error)
}
