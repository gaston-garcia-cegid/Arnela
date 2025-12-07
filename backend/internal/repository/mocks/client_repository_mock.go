package mocks

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockClientRepository is a mock implementation of ClientRepository
type MockClientRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockClientRepository) Create(ctx context.Context, client *domain.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

// GetByID mocks the GetByID method
func (m *MockClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

// GetByEmail mocks the GetByEmail method
func (m *MockClientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

// GetByDNICIF mocks the GetByDNICIF method
func (m *MockClientRepository) GetByDNICIF(ctx context.Context, dnicif string) (*domain.Client, error) {
	args := m.Called(ctx, dnicif)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

// GetByUserID mocks the GetByUserID method
func (m *MockClientRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

// Update mocks the Update method
func (m *MockClientRepository) Update(ctx context.Context, client *domain.Client) error {
	args := m.Called(ctx, client)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks the List method
func (m *MockClientRepository) List(ctx context.Context, filters repository.ClientFilters, offset, limit int) ([]*domain.Client, error) {
	args := m.Called(ctx, filters, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Client), args.Error(1)
}

// Count mocks the Count method
func (m *MockClientRepository) Count(ctx context.Context, filters repository.ClientFilters) (int, error) {
	args := m.Called(ctx, filters)
	return args.Int(0), args.Error(1)
}

// EmailExists mocks the EmailExists method
func (m *MockClientRepository) EmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, email, excludeID)
	return args.Bool(0), args.Error(1)
}

// DNICIFExists mocks the DNICIFExists method
func (m *MockClientRepository) DNICIFExists(ctx context.Context, dnicif string, excludeID *uuid.UUID) (bool, error) {
	args := m.Called(ctx, dnicif, excludeID)
	return args.Bool(0), args.Error(1)
}

// FindDeletedByEmailOrDNI mocks the FindDeletedByEmailOrDNI method
func (m *MockClientRepository) FindDeletedByEmailOrDNI(ctx context.Context, email, dnicif string) (*domain.Client, error) {
	args := m.Called(ctx, email, dnicif)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Client), args.Error(1)
}

// Reactivate mocks the Reactivate method
func (m *MockClientRepository) Reactivate(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
