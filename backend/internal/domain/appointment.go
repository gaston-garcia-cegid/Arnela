package domain

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AppointmentStatus represents the status of an appointment
type AppointmentStatus string

const (
	AppointmentStatusPending     AppointmentStatus = "pending"
	AppointmentStatusConfirmed   AppointmentStatus = "confirmed"
	AppointmentStatusCancelled   AppointmentStatus = "cancelled"
	AppointmentStatusCompleted   AppointmentStatus = "completed"
	AppointmentStatusRescheduled AppointmentStatus = "rescheduled"
)

// NullableString wraps sql.NullString for custom JSON marshaling
type NullableString struct {
	sql.NullString
}

// MarshalJSON customizes JSON output for NullableString
func (ns NullableString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return json.Marshal(nil)
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON customizes JSON input for NullableString
func (ns *NullableString) UnmarshalJSON(data []byte) error {
	var str *string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	if str == nil {
		ns.Valid = false
		ns.String = ""
	} else {
		ns.Valid = true
		ns.String = *str
	}

	return nil
}

// Appointment represents an appointment between a client and therapist
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
	Notes                 NullableString    `json:"notes" db:"notes"`                                    // ✅ Custom type
	CancellationReason    NullableString    `json:"cancellationReason" db:"cancellation_reason"`         // ✅ Custom type
	GoogleCalendarEventID NullableString    `json:"googleCalendarEventId" db:"google_calendar_event_id"` // ✅ Custom type
	CreatedBy             uuid.UUID         `json:"createdBy" db:"created_by"`
	CreatedAt             time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time         `json:"updatedAt" db:"updated_at"`
	DeletedAt             sql.NullTime      `json:"deletedAt,omitempty" db:"deleted_at"`

	// Relations (not in DB)
	Therapist *Therapist `json:"therapist,omitempty" db:"-"`
	Client    *Client    `json:"client,omitempty" db:"-"`
}

// IsDuringBusinessHours checks if appointment is Monday-Friday 9:00-18:00
func (a *Appointment) IsDuringBusinessHours() bool {
	weekday := a.StartTime.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	hour := a.StartTime.Hour()
	endHour := a.EndTime.Hour()
	endMinute := a.EndTime.Minute()

	if hour < 9 || endHour > 18 || (endHour == 18 && endMinute > 0) {
		return false
	}

	return true
}

func (a *Appointment) IsEditable() bool {
	if a.Status == AppointmentStatusCancelled || a.Status == AppointmentStatusCompleted {
		return false
	}
	return a.StartTime.After(time.Now())
}

func (a *Appointment) CanBeCancelledByClient() bool {
	if a.Status == AppointmentStatusCancelled || a.Status == AppointmentStatusCompleted {
		return false
	}
	return a.StartTime.After(time.Now().Add(24 * time.Hour))
}

// CreateAppointmentRequest represents the request to create an appointment
type CreateAppointmentRequest struct {
	ClientID        string    `json:"clientId"` // Optional: For admin/employee creating appointments for others
	TherapistID     string    `json:"therapistId" binding:"required"`
	Title           string    `json:"title" binding:"required"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"startTime" binding:"required"`
	DurationMinutes int       `json:"durationMinutes" binding:"required,oneof=45 60"`
}

// UpdateAppointmentRequest represents the request to update an appointment
type UpdateAppointmentRequest struct {
	TherapistID     string    `json:"therapistId"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	StartTime       time.Time `json:"startTime"`
	DurationMinutes int       `json:"durationMinutes"`
}

// CancelAppointmentRequest represents the request to cancel an appointment
type CancelAppointmentRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ConfirmAppointmentRequest represents the request to confirm an appointment
type ConfirmAppointmentRequest struct {
	Notes string `json:"notes"`
}

// AppointmentFilter represents filters for listing appointments
type AppointmentFilter struct {
	ClientID    *uuid.UUID         `json:"clientId"`
	TherapistID *string            `json:"therapistId"`
	Status      *AppointmentStatus `json:"status"`
	StartDate   *time.Time         `json:"startDate"`
	EndDate     *time.Time         `json:"endDate"`
	Page        int                `json:"page"`
	PageSize    int                `json:"pageSize"`
}
