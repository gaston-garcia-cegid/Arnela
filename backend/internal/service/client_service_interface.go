package service

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
)

// ClientServiceInterface defines the interface for client operations
type ClientServiceInterface interface {
	CreateClient(ctx context.Context, req CreateClientRequest) (*domain.Client, error)
	GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error)
	GetClientByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error)
	UpdateClient(ctx context.Context, id uuid.UUID, req UpdateClientRequest) (*domain.Client, error)
	DeleteClient(ctx context.Context, id uuid.UUID) error
	ListClients(ctx context.Context, filters repository.ClientFilters, page, pageSize int) (*ClientListResponse, error)
}

// CreateClientRequest represents a request to create a new client
type CreateClientRequest struct {
	UserID      *uuid.UUID `json:"userId,omitempty"`
	FirstName   string     `json:"firstName" binding:"required"`
	LastName    string     `json:"lastName" binding:"required"`
	Email       string     `json:"email" binding:"required,email"`
	Phone       string     `json:"phone" binding:"required"`
	NIF         string     `json:"nif" binding:"required"`
	DNI         string     `json:"dni" binding:"required"`
	DateOfBirth *string    `json:"dateOfBirth,omitempty"` // Format: YYYY-MM-DD
	Address     string     `json:"address,omitempty"`
	City        string     `json:"city,omitempty"`
	PostalCode  string     `json:"postalCode,omitempty"`
	Province    string     `json:"province,omitempty"`
	Notes       string     `json:"notes,omitempty"`
}

// UpdateClientRequest represents a request to update a client
type UpdateClientRequest struct {
	UserID      *uuid.UUID `json:"userId,omitempty"`
	FirstName   *string    `json:"firstName,omitempty"`
	LastName    *string    `json:"lastName,omitempty"`
	Email       *string    `json:"email,omitempty" binding:"omitempty,email"`
	Phone       *string    `json:"phone,omitempty"`
	NIF         *string    `json:"nif,omitempty"`
	DNI         *string    `json:"dni,omitempty"`
	DateOfBirth *string    `json:"dateOfBirth,omitempty"` // Format: YYYY-MM-DD
	Address     *string    `json:"address,omitempty"`
	City        *string    `json:"city,omitempty"`
	PostalCode  *string    `json:"postalCode,omitempty"`
	Province    *string    `json:"province,omitempty"`
	IsActive    *bool      `json:"isActive,omitempty"`
	Notes       *string    `json:"notes,omitempty"`
}

// ClientListResponse represents a paginated list of clients
type ClientListResponse struct {
	Clients    []*domain.Client `json:"clients"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalPages int              `json:"totalPages"`
}
