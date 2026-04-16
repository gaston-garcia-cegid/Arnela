package gcal

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// EventSync abstracts Google Calendar operations for tests and production.
type EventSync interface {
	CreateEvent(ctx context.Context, appt *domain.Appointment) (string, error)
	UpdateEvent(ctx context.Context, eventID string, appt *domain.Appointment) error
	DeleteEvent(ctx context.Context, eventID string) error
}

// AppointmentCalendarRepository loads and persists calendar linkage for appointments.
type AppointmentCalendarRepository interface {
	GetByIDWithRelations(ctx context.Context, id uuid.UUID) (*domain.Appointment, error)
	Update(ctx context.Context, appointment *domain.Appointment) error
}
