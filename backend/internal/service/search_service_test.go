package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchRepository is a mock implementation of the search repository
type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) SearchClients(ctx context.Context, query string, limit int) ([]domain.SearchClient, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.SearchClient), args.Error(1)
}

func (m *MockSearchRepository) SearchEmployees(ctx context.Context, query string, limit int) ([]domain.SearchEmployee, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.SearchEmployee), args.Error(1)
}

func (m *MockSearchRepository) SearchAppointments(ctx context.Context, query string, limit int) ([]domain.SearchAppointment, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.SearchAppointment), args.Error(1)
}

func (m *MockSearchRepository) SearchInvoices(ctx context.Context, query string, limit int) ([]domain.SearchInvoice, error) {
	args := m.Called(ctx, query, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.SearchInvoice), args.Error(1)
}

func TestSearchService_GlobalSearch_Success(t *testing.T) {
	mockRepo := new(MockSearchRepository)
	service := NewSearchService(mockRepo)

	query := "test"
	limit := 5

	// Mock data
	clientID := uuid.New()
	employeeID := uuid.New()
	appointmentID := uuid.New()
	invoiceID := uuid.New()

	mockClients := []domain.SearchClient{
		{
			ID:        clientID,
			FirstName: "Test",
			LastName:  "Client",
			Email:     "test@example.com",
		},
	}

	mockEmployees := []domain.SearchEmployee{
		{
			ID:    employeeID,
			Name:  "Test Employee",
			Email: "employee@test.com",
		},
	}

	mockAppointments := []domain.SearchAppointment{
		{
			ID:         appointmentID,
			Title:      "Test Appointment",
			StartTime:  time.Now(),
			EndTime:    time.Now().Add(1 * time.Hour),
			Status:     "confirmed",
			ClientName: "Test Client",
		},
	}

	mockInvoices := []domain.SearchInvoice{
		{
			ID:            invoiceID,
			InvoiceNumber: "TEST-001",
			ClientName:    "Test Client",
			TotalAmount:   100.00,
			Status:        "paid",
		},
	}

	ctx := context.Background()
	mockRepo.On("SearchClients", ctx, strings.ToLower(query), limit).Return(mockClients, nil)
	mockRepo.On("SearchEmployees", ctx, strings.ToLower(query), limit).Return(mockEmployees, nil)
	mockRepo.On("SearchAppointments", ctx, strings.ToLower(query), limit).Return(mockAppointments, nil)
	mockRepo.On("SearchInvoices", ctx, strings.ToLower(query), limit).Return(mockInvoices, nil)

	// Execute
	results, err := service.GlobalSearch(query)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 4, results.TotalResults)
	assert.Len(t, results.Clients, 1)
	assert.Len(t, results.Employees, 1)
	assert.Len(t, results.Appointments, 1)
	assert.Len(t, results.Invoices, 1)

	mockRepo.AssertExpectations(t)
}

func TestSearchService_GlobalSearch_NoResults(t *testing.T) {
	mockRepo := new(MockSearchRepository)
	service := NewSearchService(mockRepo)

	query := "nonexistent"
	limit := 5

	ctx := context.Background()
	mockRepo.On("SearchClients", ctx, strings.ToLower(query), limit).Return([]domain.SearchClient{}, nil)
	mockRepo.On("SearchEmployees", ctx, strings.ToLower(query), limit).Return([]domain.SearchEmployee{}, nil)
	mockRepo.On("SearchAppointments", ctx, strings.ToLower(query), limit).Return([]domain.SearchAppointment{}, nil)
	mockRepo.On("SearchInvoices", ctx, strings.ToLower(query), limit).Return([]domain.SearchInvoice{}, nil)

	// Execute
	results, err := service.GlobalSearch(query)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Equal(t, 0, results.TotalResults)
	assert.Empty(t, results.Clients)
	assert.Empty(t, results.Employees)
	assert.Empty(t, results.Appointments)
	assert.Empty(t, results.Invoices)

	mockRepo.AssertExpectations(t)
}

func TestSearchService_GlobalSearch_PartialError(t *testing.T) {
	mockRepo := new(MockSearchRepository)
	service := NewSearchService(mockRepo)

	query := "test"
	limit := 5

	mockClients := []domain.SearchClient{
		{
			ID:        uuid.New(),
			FirstName: "Test",
			LastName:  "Client",
		},
	}

	ctx := context.Background()
	mockRepo.On("SearchClients", ctx, strings.ToLower(query), limit).Return(mockClients, nil)
	mockRepo.On("SearchEmployees", ctx, strings.ToLower(query), limit).Return(nil, assert.AnError)
	mockRepo.On("SearchAppointments", ctx, strings.ToLower(query), limit).Return([]domain.SearchAppointment{}, nil)
	mockRepo.On("SearchInvoices", ctx, strings.ToLower(query), limit).Return([]domain.SearchInvoice{}, nil)

	// Execute - should still return results from successful searches
	results, err := service.GlobalSearch(query)

	// Assert - partial results should still be returned
	assert.NoError(t, err, "Service should handle partial errors gracefully")
	assert.NotNil(t, results)
	assert.Equal(t, 1, results.TotalResults)
	assert.Len(t, results.Clients, 1)
	assert.Empty(t, results.Employees) // This one failed but shouldn't break the search

	mockRepo.AssertExpectations(t)
}

func TestSearchService_GlobalSearch_CaseInsensitive(t *testing.T) {
	mockRepo := new(MockSearchRepository)
	service := NewSearchService(mockRepo)

	query := "TEST"
	expectedLowerQuery := strings.ToLower(query)
	limit := 5

	ctx := context.Background()
	mockRepo.On("SearchClients", ctx, expectedLowerQuery, limit).Return([]domain.SearchClient{}, nil)
	mockRepo.On("SearchEmployees", ctx, expectedLowerQuery, limit).Return([]domain.SearchEmployee{}, nil)
	mockRepo.On("SearchAppointments", ctx, expectedLowerQuery, limit).Return([]domain.SearchAppointment{}, nil)
	mockRepo.On("SearchInvoices", ctx, expectedLowerQuery, limit).Return([]domain.SearchInvoice{}, nil)

	// Execute
	results, err := service.GlobalSearch(query)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, results)

	mockRepo.AssertExpectations(t)
}

func TestSearchService_GlobalSearch_LimitPerType(t *testing.T) {
	mockRepo := new(MockSearchRepository)
	service := NewSearchService(mockRepo)

	query := "test"
	limit := 5

	// Create exactly 5 results (max limit)
	mockClients := make([]domain.SearchClient, 5)
	for i := 0; i < 5; i++ {
		mockClients[i] = domain.SearchClient{
			ID:        uuid.New(),
			FirstName: "Test",
			LastName:  string(rune('A' + i)),
		}
	}

	ctx := context.Background()
	mockRepo.On("SearchClients", ctx, strings.ToLower(query), limit).Return(mockClients, nil)
	mockRepo.On("SearchEmployees", ctx, strings.ToLower(query), limit).Return([]domain.SearchEmployee{}, nil)
	mockRepo.On("SearchAppointments", ctx, strings.ToLower(query), limit).Return([]domain.SearchAppointment{}, nil)
	mockRepo.On("SearchInvoices", ctx, strings.ToLower(query), limit).Return([]domain.SearchInvoice{}, nil)

	// Execute
	results, err := service.GlobalSearch(query)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results.Clients, 5, "Should respect limit of 5 per type")

	mockRepo.AssertExpectations(t)
}
