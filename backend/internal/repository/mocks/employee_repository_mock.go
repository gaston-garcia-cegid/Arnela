package mocks

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeRepository is a mock implementation of EmployeeRepository
type MockEmployeeRepository struct {
	mock.Mock
}

func (m *MockEmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) GetByDNI(ctx context.Context, dni string) (*domain.Employee, error) {
	args := m.Called(ctx, dni)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	args := m.Called(ctx, employee)
	return args.Error(0)
}

func (m *MockEmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmployeeRepository) List(ctx context.Context, limit, offset int) ([]*domain.Employee, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockEmployeeRepository) GetBySpecialty(ctx context.Context, specialty string) ([]*domain.Employee, error) {
	args := m.Called(ctx, specialty)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Employee), args.Error(1)
}

func (m *MockEmployeeRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockEmployeeRepository) DNIExists(ctx context.Context, dni string) (bool, error) {
	args := m.Called(ctx, dni)
	return args.Bool(0), args.Error(1)
}
