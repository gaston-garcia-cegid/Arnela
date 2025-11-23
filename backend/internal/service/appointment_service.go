package service

import (
	"context"
	"fmt"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
)

// AppointmentServiceInterface defines the interface for appointment service
type AppointmentServiceInterface interface {
	// Client operations
	CreateAppointment(ctx context.Context, req domain.CreateAppointmentRequest, createdBy uuid.UUID) (*domain.Appointment, error)
	GetAppointment(ctx context.Context, id uuid.UUID) (*domain.Appointment, error)
	UpdateAppointment(ctx context.Context, id uuid.UUID, req domain.UpdateAppointmentRequest, userID uuid.UUID) (*domain.Appointment, error)
	CancelAppointment(ctx context.Context, id uuid.UUID, req domain.CancelAppointmentRequest, userID uuid.UUID, isAdmin bool) error
	GetMyAppointments(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*domain.Appointment, int, error)

	// Admin operations
	ConfirmAppointment(ctx context.Context, id uuid.UUID, req domain.ConfirmAppointmentRequest) (*domain.Appointment, error)
	ListAppointments(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, int, error)
	GetAppointmentsByTherapist(ctx context.Context, therapistID string, startDate, endDate time.Time) ([]*domain.Appointment, error)
	GetAvailableSlots(ctx context.Context, therapistID string, date time.Time, duration int) ([]time.Time, error)

	// Utility
	GetTherapists(ctx context.Context) []domain.Therapist
	ValidateAppointmentTime(ctx context.Context, therapistID string, startTime time.Time, duration int, excludeID *uuid.UUID) error
}

type appointmentService struct {
	appointmentRepo repository.AppointmentRepository
	clientRepo      repository.ClientRepository
}

func NewAppointmentService(appointmentRepo repository.AppointmentRepository, clientRepo repository.ClientRepository) AppointmentServiceInterface {
	return &appointmentService{
		appointmentRepo: appointmentRepo,
		clientRepo:      clientRepo,
	}
}

// CreateAppointment creates a new appointment with pending status
func (s *appointmentService) CreateAppointment(ctx context.Context, req domain.CreateAppointmentRequest, createdBy uuid.UUID) (*domain.Appointment, error) {
	// Validate client exists
	client, err := s.clientRepo.GetByID(ctx, req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("cliente no encontrado")
	}

	if !client.IsActive {
		return nil, fmt.Errorf("el cliente está inactivo")
	}

	// Validate therapist exists
	if !domain.IsValidTherapistID(req.TherapistID) {
		return nil, fmt.Errorf("terapeuta no válido")
	}

	// Validate duration
	if req.DurationMinutes != 45 && req.DurationMinutes != 60 {
		return nil, fmt.Errorf("la duración debe ser 45 o 60 minutos")
	}

	// Calculate end time
	endTime := req.StartTime.Add(time.Duration(req.DurationMinutes) * time.Minute)

	// Create appointment object for validation
	appointment := &domain.Appointment{
		StartTime:       req.StartTime,
		EndTime:         endTime,
		DurationMinutes: req.DurationMinutes,
		TherapistID:     req.TherapistID,
	}

	// Validate business hours (Mon-Fri 9:00-18:00)
	if !appointment.IsDuringBusinessHours() {
		return nil, fmt.Errorf("la cita debe ser de lunes a viernes entre las 9:00 y las 18:00")
	}

	// Validate time is in the future
	if req.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("la cita debe ser en el futuro")
	}

	// Validate no overlap (with 15min buffer)
	if err := s.ValidateAppointmentTime(ctx, req.TherapistID, req.StartTime, req.DurationMinutes, nil); err != nil {
		return nil, err
	}

	// Create appointment
	appointment.ID = uuid.New()
	appointment.ClientID = req.ClientID
	appointment.Title = req.Title
	appointment.Description = req.Description
	appointment.Status = domain.AppointmentStatusPending
	appointment.CreatedBy = createdBy
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()

	if err := s.appointmentRepo.Create(ctx, appointment); err != nil {
		return nil, fmt.Errorf("failed to create appointment: %w", err)
	}

	// Load relations for response
	return s.appointmentRepo.GetByIDWithRelations(ctx, appointment.ID)
}

// GetAppointment retrieves an appointment by ID
func (s *appointmentService) GetAppointment(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	return s.appointmentRepo.GetByIDWithRelations(ctx, id)
}

// UpdateAppointment updates an appointment (only if editable by client)
func (s *appointmentService) UpdateAppointment(ctx context.Context, id uuid.UUID, req domain.UpdateAppointmentRequest, userID uuid.UUID) (*domain.Appointment, error) {
	// Get existing appointment
	appointment, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cita no encontrada")
	}

	// Check if appointment is editable
	if !appointment.IsEditable() {
		return nil, fmt.Errorf("la cita no puede ser modificada (ya pasó o está cancelada)")
	}

	// Validate client owns the appointment
	client, err := s.clientRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("cliente no encontrado")
	}

	if appointment.ClientID != client.ID {
		return nil, fmt.Errorf("no tienes permiso para modificar esta cita")
	}

	// Update fields
	if req.Title != "" {
		appointment.Title = req.Title
	}

	if req.Description != "" {
		appointment.Description = req.Description
	}

	if req.TherapistID != "" {
		if !domain.IsValidTherapistID(req.TherapistID) {
			return nil, fmt.Errorf("terapeuta no válido")
		}
		appointment.TherapistID = req.TherapistID
	}

	// Handle time updates
	if !req.StartTime.IsZero() || req.DurationMinutes > 0 {
		newStartTime := appointment.StartTime
		newDuration := appointment.DurationMinutes

		if !req.StartTime.IsZero() {
			newStartTime = req.StartTime

			// Validate time is in the future
			if newStartTime.Before(time.Now()) {
				return nil, fmt.Errorf("la cita debe ser en el futuro")
			}
		}

		if req.DurationMinutes > 0 {
			if req.DurationMinutes != 45 && req.DurationMinutes != 60 {
				return nil, fmt.Errorf("la duración debe ser 45 o 60 minutos")
			}
			newDuration = req.DurationMinutes
		}

		newEndTime := newStartTime.Add(time.Duration(newDuration) * time.Minute)

		// Create temp appointment for validation
		tempAppt := &domain.Appointment{
			StartTime:       newStartTime,
			EndTime:         newEndTime,
			DurationMinutes: newDuration,
		}

		// Validate business hours
		if !tempAppt.IsDuringBusinessHours() {
			return nil, fmt.Errorf("la cita debe ser de lunes a viernes entre las 9:00 y las 18:00")
		}

		// Validate no overlap (excluding current appointment)
		if err := s.ValidateAppointmentTime(ctx, appointment.TherapistID, newStartTime, newDuration, &id); err != nil {
			return nil, err
		}

		appointment.StartTime = newStartTime
		appointment.EndTime = newEndTime
		appointment.DurationMinutes = newDuration

		// Mark as rescheduled if time changed
		if appointment.Status == domain.AppointmentStatusConfirmed {
			appointment.Status = domain.AppointmentStatusRescheduled
		}
	}

	appointment.UpdatedAt = time.Now()

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return s.appointmentRepo.GetByIDWithRelations(ctx, id)
}

// CancelAppointment cancels an appointment
func (s *appointmentService) CancelAppointment(ctx context.Context, id uuid.UUID, req domain.CancelAppointmentRequest, userID uuid.UUID, isAdmin bool) error {
	appointment, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("cita no encontrada")
	}

	// If not admin, check if client owns the appointment and can cancel
	if !isAdmin {
		client, err := s.clientRepo.GetByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("cliente no encontrado")
		}

		if appointment.ClientID != client.ID {
			return fmt.Errorf("no tienes permiso para cancelar esta cita")
		}

		if !appointment.CanBeCancelledByClient() {
			return fmt.Errorf("la cita no puede ser cancelada (ya pasó o ya está cancelada)")
		}
	}

	// Update appointment
	appointment.Status = domain.AppointmentStatusCancelled
	appointment.CancellationReason = req.Reason
	appointment.UpdatedAt = time.Now()

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		return fmt.Errorf("failed to cancel appointment: %w", err)
	}

	return nil
}

// GetMyAppointments retrieves appointments for a specific client
func (s *appointmentService) GetMyAppointments(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*domain.Appointment, int, error) {
	filters := domain.AppointmentFilter{
		ClientID: &clientID,
		Page:     page,
		PageSize: pageSize,
	}

	appointments, err := s.appointmentRepo.ListWithRelations(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get appointments: %w", err)
	}

	count, err := s.appointmentRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count appointments: %w", err)
	}

	return appointments, count, nil
}

// ConfirmAppointment confirms a pending appointment (admin only)
func (s *appointmentService) ConfirmAppointment(ctx context.Context, id uuid.UUID, req domain.ConfirmAppointmentRequest) (*domain.Appointment, error) {
	appointment, err := s.appointmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("cita no encontrada")
	}

	if appointment.Status != domain.AppointmentStatusPending {
		return nil, fmt.Errorf("solo se pueden confirmar citas pendientes")
	}

	appointment.Status = domain.AppointmentStatusConfirmed
	appointment.Notes = req.Notes
	appointment.UpdatedAt = time.Now()

	if err := s.appointmentRepo.Update(ctx, appointment); err != nil {
		return nil, fmt.Errorf("failed to confirm appointment: %w", err)
	}

	return s.appointmentRepo.GetByIDWithRelations(ctx, id)
}

// ListAppointments lists all appointments with filters (admin only)
func (s *appointmentService) ListAppointments(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, int, error) {
	appointments, err := s.appointmentRepo.ListWithRelations(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list appointments: %w", err)
	}

	count, err := s.appointmentRepo.Count(ctx, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count appointments: %w", err)
	}

	return appointments, count, nil
}

// GetAppointmentsByTherapist retrieves appointments for a therapist in a date range
func (s *appointmentService) GetAppointmentsByTherapist(ctx context.Context, therapistID string, startDate, endDate time.Time) ([]*domain.Appointment, error) {
	if !domain.IsValidTherapistID(therapistID) {
		return nil, fmt.Errorf("terapeuta no válido")
	}

	return s.appointmentRepo.GetByDateRange(ctx, startDate, endDate, &therapistID)
}

// GetAvailableSlots returns available time slots for a therapist on a specific date
func (s *appointmentService) GetAvailableSlots(ctx context.Context, therapistID string, date time.Time, duration int) ([]time.Time, error) {
	if !domain.IsValidTherapistID(therapistID) {
		return nil, fmt.Errorf("terapeuta no válido")
	}

	if duration != 45 && duration != 60 {
		return nil, fmt.Errorf("la duración debe ser 45 o 60 minutos")
	}

	// Check if date is a weekday
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return []time.Time{}, nil // No slots on weekends
	}

	// Get start and end of business day
	year, month, day := date.Date()
	startOfDay := time.Date(year, month, day, 9, 0, 0, 0, date.Location())
	endOfDay := time.Date(year, month, day, 18, 0, 0, 0, date.Location())

	// Get existing appointments for the day
	existingAppointments, err := s.appointmentRepo.GetByDateRange(ctx, startOfDay, endOfDay, &therapistID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing appointments: %w", err)
	}

	// Generate all possible slots (15min intervals)
	var availableSlots []time.Time
	currentSlot := startOfDay

	for currentSlot.Before(endOfDay) {
		slotEndTime := currentSlot.Add(time.Duration(duration) * time.Minute)

		// Check if slot would exceed business hours
		if slotEndTime.After(endOfDay) {
			break
		}

		// Check if slot overlaps with any existing appointment (with 15min buffer)
		isAvailable := true
		buffer := 15 * time.Minute

		for _, appt := range existingAppointments {
			// Skip cancelled appointments
			if appt.Status == domain.AppointmentStatusCancelled {
				continue
			}

			// Check overlap with buffer
			if currentSlot.Add(-buffer).Before(appt.EndTime) && slotEndTime.Add(buffer).After(appt.StartTime) {
				isAvailable = false
				break
			}
		}

		if isAvailable {
			availableSlots = append(availableSlots, currentSlot)
		}

		// Move to next slot (15min interval)
		currentSlot = currentSlot.Add(15 * time.Minute)
	}

	return availableSlots, nil
}

// GetTherapists returns all available therapists
func (s *appointmentService) GetTherapists(ctx context.Context) []domain.Therapist {
	return domain.GetMockTherapists()
}

// ValidateAppointmentTime validates if an appointment time is valid (no overlap, business hours)
func (s *appointmentService) ValidateAppointmentTime(ctx context.Context, therapistID string, startTime time.Time, duration int, excludeID *uuid.UUID) error {
	endTime := startTime.Add(time.Duration(duration) * time.Minute)

	// Add 15min buffer before and after
	bufferStartTime := startTime.Add(-15 * time.Minute)
	bufferEndTime := endTime.Add(15 * time.Minute)

	// Check for overlapping appointments
	hasOverlap, err := s.appointmentRepo.CheckOverlap(ctx, therapistID, bufferStartTime, bufferEndTime, excludeID)
	if err != nil {
		return fmt.Errorf("failed to check overlap: %w", err)
	}

	if hasOverlap {
		return fmt.Errorf("el horario no está disponible (conflicto con otra cita)")
	}

	return nil
}
