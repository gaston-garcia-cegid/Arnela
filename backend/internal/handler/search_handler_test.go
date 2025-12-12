package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchService is a mock implementation of the search service
type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) GlobalSearch(query string) (*domain.SearchResults, error) {
	args := m.Called(query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SearchResults), args.Error(1)
}

func setupSearchTestRouter(mockService *MockSearchService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := NewSearchHandler(mockService)

	// Protected route (requires auth middleware in production)
	api := r.Group("/api/v1")
	{
		api.GET("/search", handler.GlobalSearch)
	}

	return r
}

func TestGlobalSearch_Success(t *testing.T) {
	mockService := new(MockSearchService)
	router := setupSearchTestRouter(mockService)

	// Mock data
	clientID := uuid.New()
	employeeID := uuid.New()
	appointmentID := uuid.New()
	invoiceID := uuid.New()

	expectedResults := &domain.SearchResults{
		Clients: []domain.SearchClient{
			{
				ID:        clientID,
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
				Phone:     "+34666777888",
				DNICIF:    "12345678Z",
			},
		},
		Employees: []domain.SearchEmployee{
			{
				ID:          employeeID,
				Name:        "Dr. María García",
				Email:       "maria@clinic.com",
				Phone:       "+34600111222",
				Specialties: []string{"Fisioterapia", "Osteopatía"},
				AvatarColor: "#FF5733",
			},
		},
		Appointments: []domain.SearchAppointment{
			{
				ID:           appointmentID,
				Title:        "Consulta Fisioterapia",
				StartTime:    time.Now(),
				EndTime:      time.Now().Add(1 * time.Hour),
				Status:       "confirmed",
				ClientName:   "Juan Pérez",
				EmployeeName: "Dr. María García",
			},
		},
		Invoices: []domain.SearchInvoice{
			{
				ID:            invoiceID,
				InvoiceNumber: "FAC-2025-001",
				ClientName:    "Juan Pérez",
				TotalAmount:   50.00,
				Status:        "paid",
				IssueDate:     time.Now(),
			},
		},
		TotalResults: 4,
	}

	mockService.On("GlobalSearch", "juan").Return(expectedResults, nil)

	// Make request
	req, _ := http.NewRequest("GET", "/api/v1/search?q=juan", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.SearchResults
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 4, response.TotalResults)
	assert.Len(t, response.Clients, 1)
	assert.Len(t, response.Employees, 1)
	assert.Len(t, response.Appointments, 1)
	assert.Len(t, response.Invoices, 1)
	assert.Equal(t, "Juan", response.Clients[0].FirstName)
	assert.Equal(t, "Dr. María García", response.Employees[0].Name)

	mockService.AssertExpectations(t)
}

func TestGlobalSearch_EmptyQuery(t *testing.T) {
	mockService := new(MockSearchService)
	router := setupSearchTestRouter(mockService)

	// Make request with empty query
	req, _ := http.NewRequest("GET", "/api/v1/search?q=", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "query parameter 'q' is required")

	mockService.AssertNotCalled(t, "GlobalSearch", mock.Anything)
}

func TestGlobalSearch_QueryTooShort(t *testing.T) {
	mockService := new(MockSearchService)
	router := setupSearchTestRouter(mockService)

	// Make request with query too short (less than 2 characters)
	req, _ := http.NewRequest("GET", "/api/v1/search?q=a", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response["error"], "at least 2 characters")

	mockService.AssertNotCalled(t, "GlobalSearch", mock.Anything)
}

func TestGlobalSearch_NoResults(t *testing.T) {
	mockService := new(MockSearchService)
	router := setupSearchTestRouter(mockService)

	emptyResults := &domain.SearchResults{
		Clients:      []domain.SearchClient{},
		Employees:    []domain.SearchEmployee{},
		Appointments: []domain.SearchAppointment{},
		Invoices:     []domain.SearchInvoice{},
		TotalResults: 0,
	}

	mockService.On("GlobalSearch", "nonexistent").Return(emptyResults, nil)

	// Make request
	req, _ := http.NewRequest("GET", "/api/v1/search?q=nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.SearchResults
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 0, response.TotalResults)
	assert.Empty(t, response.Clients)
	assert.Empty(t, response.Employees)
	assert.Empty(t, response.Appointments)
	assert.Empty(t, response.Invoices)

	mockService.AssertExpectations(t)
}

func TestGlobalSearch_ServiceError(t *testing.T) {
	mockService := new(MockSearchService)
	router := setupSearchTestRouter(mockService)

	mockService.On("GlobalSearch", "error").Return(nil, assert.AnError)

	// Make request
	req, _ := http.NewRequest("GET", "/api/v1/search?q=error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["error"])

	mockService.AssertExpectations(t)
}

func TestGlobalSearch_CaseInsensitive(t *testing.T) {
	mockService := new(MockSearchService)
	router := setupSearchTestRouter(mockService)

	expectedResults := &domain.SearchResults{
		Clients: []domain.SearchClient{
			{
				ID:        uuid.New(),
				FirstName: "Juan",
				LastName:  "Pérez",
				Email:     "juan@example.com",
			},
		},
		Employees:    []domain.SearchEmployee{},
		Appointments: []domain.SearchAppointment{},
		Invoices:     []domain.SearchInvoice{},
		TotalResults: 1,
	}

	// Service should receive lowercase query
	mockService.On("GlobalSearch", "JUAN").Return(expectedResults, nil)

	// Make request with uppercase query
	req, _ := http.NewRequest("GET", "/api/v1/search?q=JUAN", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.SearchResults
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.TotalResults)
	assert.Len(t, response.Clients, 1)

	mockService.AssertExpectations(t)
}

func TestGlobalSearch_MaxResultsPerType(t *testing.T) {
	mockService := new(MockSearchService)
	router := setupSearchTestRouter(mockService)

	// Create 6 clients (should return max 5)
	clients := make([]domain.SearchClient, 5)
	for i := 0; i < 5; i++ {
		clients[i] = domain.SearchClient{
			ID:        uuid.New(),
			FirstName: "Client",
			LastName:  string(rune('A' + i)),
			Email:     "client@example.com",
		}
	}

	expectedResults := &domain.SearchResults{
		Clients:      clients,
		Employees:    []domain.SearchEmployee{},
		Appointments: []domain.SearchAppointment{},
		Invoices:     []domain.SearchInvoice{},
		TotalResults: 5,
	}

	mockService.On("GlobalSearch", "client").Return(expectedResults, nil)

	// Make request
	req, _ := http.NewRequest("GET", "/api/v1/search?q=client", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.SearchResults
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.LessOrEqual(t, len(response.Clients), 5, "Should return max 5 clients")

	mockService.AssertExpectations(t)
}
