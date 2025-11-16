package domain

import (
	"time"

	"github.com/google/uuid"
)

// Client represents a client in the system
type Client struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	UserID      *uuid.UUID `json:"userId,omitempty" db:"user_id"` // Optional: if client has login access
	FirstName   string     `json:"firstName" db:"first_name"`
	LastName    string     `json:"lastName" db:"last_name"`
	Email       string     `json:"email" db:"email"`
	Phone       string     `json:"phone" db:"phone"`
	NIF         string     `json:"nif" db:"nif"`           // Spanish Tax ID (NIF) - Required
	DNI         string     `json:"dni,omitempty" db:"dni"` // Spanish National ID (DNI/NIE) - Optional
	DateOfBirth *time.Time `json:"dateOfBirth,omitempty" db:"date_of_birth"`
	Address     string     `json:"address,omitempty" db:"address"`
	City        string     `json:"city,omitempty" db:"city"`
	PostalCode  string     `json:"postalCode,omitempty" db:"postal_code"`
	Province    string     `json:"province,omitempty" db:"province"`
	IsActive    bool       `json:"isActive" db:"is_active"`
	LastVisit   *time.Time `json:"lastVisit,omitempty" db:"last_visit"`
	Notes       string     `json:"notes,omitempty" db:"notes"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty" db:"deleted_at"` // Soft delete
}

// FullName returns the client's full name
func (c *Client) FullName() string {
	return c.FirstName + " " + c.LastName
}

// IsDeleted checks if the client has been soft deleted
func (c *Client) IsDeleted() bool {
	return c.DeletedAt != nil
}

// Age calculates the client's age if date of birth is set
func (c *Client) Age() *int {
	if c.DateOfBirth == nil {
		return nil
	}
	age := time.Now().Year() - c.DateOfBirth.Year()
	if time.Now().YearDay() < c.DateOfBirth.YearDay() {
		age--
	}
	return &age
}

// ClientStatus represents the status of a client
type ClientStatus string

const (
	ClientStatusActive   ClientStatus = "active"
	ClientStatusInactive ClientStatus = "inactive"
)
