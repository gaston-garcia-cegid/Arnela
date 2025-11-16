package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleEmployee UserRole = "employee"
	RoleClient   UserRole = "client"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password in JSON
	FirstName    string    `json:"firstName" db:"first_name"`
	LastName     string    `json:"lastName" db:"last_name"`
	Role         UserRole  `json:"role" db:"role"`
	IsActive     bool      `json:"isActive" db:"is_active"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEmployee checks if the user has employee role
func (u *User) IsEmployee() bool {
	return u.Role == RoleEmployee
}

// IsClient checks if the user has client role
func (u *User) IsClient() bool {
	return u.Role == RoleClient
}
