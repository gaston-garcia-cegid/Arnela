package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// AppointmentStatus represents the state of an appointment
type AppointmentStatus string

const (
	AppointmentStatusPending     AppointmentStatus = "pending"
	AppointmentStatusConfirmed   AppointmentStatus = "confirmed"
	AppointmentStatusCancelled   AppointmentStatus = "cancelled"
	AppointmentStatusCompleted   AppointmentStatus = "completed"
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

// Appointment represents a scheduled appointment with a therapist
type Appointment struct {
	ID                    uuid.UUID         `json:"id" db:"id"`
	ClientID              uuid.UUID         `json:"clientId" db:"client_id"`
	TherapistID           string            `json:"therapistId" db:"therapist_id"`
	Title                 string            `json:"title" db:"title"`
	Description           string            `json:"description" db:"description"`
	StartTime             time.Time         `json:"startTime" db:"start_time"`
	EndTime               time.Time         `json:"endTime" db:"end_time"`
	DurationMinutes       int               `json:"durationMinutes" db:"duration_minutes"`
	Status                AppointmentStatus `json:"status" db:"status"`
	Notes                 string            `json:"notes" db:"notes"`
	CancellationReason    string            `json:"cancellationReason,omitempty" db:"cancellation_reason"`
	GoogleCalendarEventID string            `json:"googleCalendarEventId,omitempty" db:"google_calendar_event_id"`
	CreatedBy             uuid.UUID         `json:"createdBy" db:"created_by"`
	CreatedAt             time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time         `json:"updatedAt" db:"updated_at"`
	DeletedAt             sql.NullTime      `json:"-" db:"deleted_at"`

	// Relations (populated from joins) - NO DB TAGS
	Client    *Client    `json:"client,omitempty"`
	Therapist *Therapist `json:"therapist,omitempty"`
}

// IsEditable returns true if the appointment can be edited by the client
func (a *Appointment) IsEditable() bool {
	// Can only edit pending or confirmed appointments
	if a.Status != AppointmentStatusPending && a.Status != AppointmentStatusConfirmed {
		return false
	}

	// Cannot edit if already started
	return time.Now().Before(a.StartTime)
}

// CanBeCancelledByClient returns true if the appointment can be cancelled by the client
func (a *Appointment) CanBeCancelledByClient() bool {
	// Can cancel pending, confirmed, or rescheduled appointments
	if a.Status != AppointmentStatusPending &&
		a.Status != AppointmentStatusConfirmed &&
		a.Status != AppointmentStatusRescheduled {
		return false
	}

	// Cannot cancel if already started
	return time.Now().Before(a.StartTime)
}

// IsDuringBusinessHours checks if the appointment is within Mon-Fri 9:00-18:00
func (a *Appointment) IsDuringBusinessHours() bool {
	// Check if it's a weekday (Monday-Friday)
	weekday := a.StartTime.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// Extract start and end hours
	startHour := a.StartTime.Hour()
	endHour := a.EndTime.Hour()
	endMinute := a.EndTime.Minute()

	// Must start at or after 9:00 and end at or before 18:00
	if startHour < 9 || endHour > 18 {
		return false
	}

	// If end hour is exactly 18, minutes must be 0
	if endHour == 18 && endMinute > 0 {
		return false
	}

	return true
}

// HasValidDuration checks if the duration is 45 or 60 minutes
func (a *Appointment) HasValidDuration() bool {
	return a.DurationMinutes == 45 || a.DurationMinutes == 60
}

// OverlapsWith checks if this appointment overlaps with another (including buffer)
func (a *Appointment) OverlapsWith(other *Appointment, bufferMinutes int) bool {
	buffer := time.Duration(bufferMinutes) * time.Minute

	// Add buffer to both appointments
	thisStart := a.StartTime.Add(-buffer)
	thisEnd := a.EndTime.Add(buffer)
	otherStart := other.StartTime.Add(-buffer)
	otherEnd := other.EndTime.Add(buffer)

	// Check for overlap
	return thisStart.Before(otherEnd) && otherStart.Before(thisEnd)
}

// CreateAppointmentRequest represents the request to create an appointment
type CreateAppointmentRequest struct {
	ClientID        uuid.UUID `json:"clientId" binding:"required"`
	TherapistID     string    `json:"therapistId" binding:"required"`
	Title           string    `json:"title" binding:"required,max=255"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"startTime" binding:"required"`
	DurationMinutes int       `json:"durationMinutes" binding:"required,oneof=45 60"`
}

// UpdateAppointmentRequest represents the request to update an appointment
type UpdateAppointmentRequest struct {
	Title           string    `json:"title" binding:"omitempty,max=255"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"startTime"`
	DurationMinutes int       `json:"durationMinutes" binding:"omitempty,oneof=45 60"`
	TherapistID     string    `json:"therapistId"`
}

// ConfirmAppointmentRequest represents the request to confirm an appointment (admin only)
type ConfirmAppointmentRequest struct {
	Notes string `json:"notes"`
}

// CancelAppointmentRequest represents the request to cancel an appointment
type CancelAppointmentRequest struct {
	Reason string `json:"reason" binding:"required,max=500"`
}

// AppointmentFilter represents filters for listing appointments
type AppointmentFilter struct {
	ClientID    *uuid.UUID
	TherapistID *string
	Status      *AppointmentStatus
	StartDate   *time.Time
	EndDate     *time.Time
	Page        int
	PageSize    int
}
