package domain

import (
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// StringArray is a custom type for PostgreSQL text array
type StringArray []string

// Scan implements the sql.Scanner interface
func (s *StringArray) Scan(src interface{}) error {
	if src == nil {
		*s = nil
		return nil
	}

	// Use pq.Array for proper PostgreSQL array scanning
	arr := pq.StringArray{}
	if err := arr.Scan(src); err != nil {
		return err
	}
	*s = StringArray(arr)
	return nil
}

// Value implements the driver.Valuer interface
func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	// Use pq.Array for proper PostgreSQL array formatting
	return pq.Array(s).Value()
}

// MarshalJSON implements json.Marshaler for proper JSON serialization
func (s StringArray) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("[]"), nil
	}
	return json.Marshal([]string(s))
}

// UnmarshalJSON implements json.Unmarshaler for proper JSON deserialization
func (s *StringArray) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	*s = StringArray(arr)
	return nil
}

// Employee represents an employee/therapist in the system
type Employee struct {
	ID          uuid.UUID   `json:"id" db:"id"`
	UserID      *uuid.UUID  `json:"userId,omitempty" db:"user_id"`
	FirstName   string      `json:"firstName" db:"first_name"`
	LastName    string      `json:"lastName" db:"last_name"`
	Email       string      `json:"email" db:"email"`
	Phone       string      `json:"phone" db:"phone"`
	DNI         string      `json:"dni" db:"dni"`
	DateOfBirth *time.Time  `json:"dateOfBirth,omitempty" db:"date_of_birth"`
	Position    *string     `json:"position,omitempty" db:"position"`
	Specialties StringArray `json:"specialties" db:"specialties"`
	IsActive    bool        `json:"isActive" db:"is_active"`
	HireDate    *time.Time  `json:"hireDate,omitempty" db:"hire_date"`
	Notes       *string     `json:"notes,omitempty" db:"notes"`
	AvatarColor string      `json:"avatarColor" db:"avatar_color"`
	CreatedAt   time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time   `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time  `json:"deletedAt,omitempty" db:"deleted_at"`
}

// FullName returns the employee's full name
func (e *Employee) FullName() string {
	return strings.TrimSpace(e.FirstName + " " + e.LastName)
}

// IsDeleted checks if the employee is soft-deleted
func (e *Employee) IsDeleted() bool {
	return e.DeletedAt != nil
}

// HasSpecialty checks if employee has a specific specialty
func (e *Employee) HasSpecialty(specialty string) bool {
	for _, s := range e.Specialties {
		if strings.EqualFold(s, specialty) {
			return true
		}
	}
	return false
}
