package service

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
)

// AuthServiceInterface defines the interface for authentication operations
type AuthServiceInterface interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetClientIDForUser(ctx context.Context, userID string) (string, error)
}
