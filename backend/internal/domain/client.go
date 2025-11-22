package domain

import (
	"time"

	"github.com/google/uuid"
)

// Client represents a client in the system
type Client struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"` // ✅ AÑADIDO
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Phone     string    `json:"phone" db:"phone"`
	DNI       string    `json:"dni" db:"dni"`
	NIF       string    `json:"nif" db:"nif"`
	Address   Address   `json:"address"`
	Notes     string    `json:"notes" db:"notes"`
	IsActive  bool      `json:"isActive" db:"is_active"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Address represents a client's address
type Address struct {
	Street     string `json:"street" db:"address_street"`
	City       string `json:"city" db:"address_city"`
	Province   string `json:"province" db:"address_province"`
	PostalCode string `json:"postalCode" db:"address_postal_code"`
	Country    string `json:"country" db:"address_country"`
}
