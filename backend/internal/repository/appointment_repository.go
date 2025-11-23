package repository

import (
	"context"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// AppointmentRepository defines the interface for appointment data operations
type AppointmentRepository interface {
	// Create creates a new appointment
	Create(ctx context.Context, appointment *domain.Appointment) error

	// GetByID retrieves an appointment by ID (excluding soft-deleted)
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error)

	// GetByIDWithRelations retrieves an appointment with client and therapist data
	GetByIDWithRelations(ctx context.Context, id uuid.UUID) (*domain.Appointment, error)

	// Update updates an appointment's information
	Update(ctx context.Context, appointment *domain.Appointment) error

	// Delete soft-deletes an appointment
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves a paginated list of appointments with optional filters
	List(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, error)

	// ListWithRelations retrieves appointments with client and therapist data
	ListWithRelations(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, error)

	// Count returns the total number of appointments matching the filters
	Count(ctx context.Context, filters domain.AppointmentFilter) (int, error)

	// GetByClientID retrieves all appointments for a client (excluding soft-deleted)
	GetByClientID(ctx context.Context, clientID uuid.UUID, limit, offset int) ([]*domain.Appointment, error)

	// GetByTherapistID retrieves all appointments for a therapist (excluding soft-deleted)
	GetByTherapistID(ctx context.Context, therapistID string, limit, offset int) ([]*domain.Appointment, error)

	// GetByDateRange retrieves appointments within a date range (excluding soft-deleted)
	GetByDateRange(ctx context.Context, startDate, endDate time.Time, therapistID *string) ([]*domain.Appointment, error)

	// CheckOverlap checks if there's an overlapping appointment for a therapist (excluding soft-deleted and current appointment)
	CheckOverlap(ctx context.Context, therapistID string, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error)

	// UpdateStatus updates only the status of an appointment
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AppointmentStatus) error
}
