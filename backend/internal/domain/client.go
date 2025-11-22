package domain

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Client represents a client in the system
type Client struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"userId" db:"user_id"`
	Email     string    `json:"email" db:"email"`
	FirstName string    `json:"firstName" db:"first_name"`
	LastName  string    `json:"lastName" db:"last_name"`
	Phone     string    `json:"phone" db:"phone"`
	DNI       string    `json:"dni" db:"dni"`
	NIF       string    `json:"nif" db:"nif"`

	// ✅ Address fields flattened with explicit db tags
	AddressStreet     string `json:"-" db:"address_street"`
	AddressCity       string `json:"-" db:"address_city"`
	AddressProvince   string `json:"-" db:"address_province"`
	AddressPostalCode string `json:"-" db:"address_postal_code"`
	AddressCountry    string `json:"-" db:"address_country"`

	Notes     string       `json:"notes" db:"notes"`
	IsActive  bool         `json:"isActive" db:"is_active"`
	CreatedAt time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time    `json:"updatedAt" db:"updated_at"`
	DeletedAt sql.NullTime `json:"-" db:"deleted_at"`
}

// Address is a computed property for JSON responses
func (c *Client) Address() Address {
	return Address{
		Street:     c.AddressStreet,
		City:       c.AddressCity,
		Province:   c.AddressProvince,
		PostalCode: c.AddressPostalCode,
		Country:    c.AddressCountry,
	}
}

// SetAddress updates address fields from Address struct
func (c *Client) SetAddress(addr Address) {
	c.AddressStreet = addr.Street
	c.AddressCity = addr.City
	c.AddressProvince = addr.Province
	c.AddressPostalCode = addr.PostalCode
	c.AddressCountry = addr.Country
}

// MarshalJSON customizes JSON output to include nested Address
func (c *Client) MarshalJSON() ([]byte, error) {
	type Alias Client
	return json.Marshal(&struct {
		*Alias
		Address Address `json:"address"`
	}{
		Alias:   (*Alias)(c),
		Address: c.Address(),
	})
}

// Address represents a client's address (for API responses)
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}
