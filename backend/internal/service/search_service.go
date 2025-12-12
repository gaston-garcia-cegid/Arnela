package service

import (
	"context"
	"strings"
	"sync"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
)

// SearchRepository defines the interface for search repository operations
type SearchRepository interface {
	SearchClients(ctx context.Context, query string, limit int) ([]domain.SearchClient, error)
	SearchEmployees(ctx context.Context, query string, limit int) ([]domain.SearchEmployee, error)
	SearchAppointments(ctx context.Context, query string, limit int) ([]domain.SearchAppointment, error)
	SearchInvoices(ctx context.Context, query string, limit int) ([]domain.SearchInvoice, error)
}

// SearchService implements domain.SearchService
type SearchService struct {
	repo SearchRepository
}

// NewSearchService creates a new search service
func NewSearchService(repo SearchRepository) *SearchService {
	return &SearchService{
		repo: repo,
	}
}

// GlobalSearch performs a search across all entity types
// It runs searches in parallel for better performance and returns aggregated results
// Partial failures are handled gracefully - if one search fails, others continue
func (s *SearchService) GlobalSearch(query string) (*domain.SearchResults, error) {
	// Normalize query to lowercase for case-insensitive search
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))

	const maxResultsPerType = 5
	ctx := context.Background()

	// Use WaitGroup to run searches in parallel
	var wg sync.WaitGroup
	var mu sync.Mutex

	results := &domain.SearchResults{
		Clients:      []domain.SearchClient{},
		Employees:    []domain.SearchEmployee{},
		Appointments: []domain.SearchAppointment{},
		Invoices:     []domain.SearchInvoice{},
		TotalResults: 0,
	}

	// Search clients
	wg.Add(1)
	go func() {
		defer wg.Done()
		clients, err := s.repo.SearchClients(ctx, normalizedQuery, maxResultsPerType)
		if err == nil && clients != nil {
			mu.Lock()
			results.Clients = clients
			results.TotalResults += len(clients)
			mu.Unlock()
		}
		// Silently handle errors - partial results are acceptable
	}()

	// Search employees
	wg.Add(1)
	go func() {
		defer wg.Done()
		employees, err := s.repo.SearchEmployees(ctx, normalizedQuery, maxResultsPerType)
		if err == nil && employees != nil {
			mu.Lock()
			results.Employees = employees
			results.TotalResults += len(employees)
			mu.Unlock()
		}
	}()

	// Search appointments
	wg.Add(1)
	go func() {
		defer wg.Done()
		appointments, err := s.repo.SearchAppointments(ctx, normalizedQuery, maxResultsPerType)
		if err == nil && appointments != nil {
			mu.Lock()
			results.Appointments = appointments
			results.TotalResults += len(appointments)
			mu.Unlock()
		}
	}()

	// Search invoices
	wg.Add(1)
	go func() {
		defer wg.Done()
		invoices, err := s.repo.SearchInvoices(ctx, normalizedQuery, maxResultsPerType)
		if err == nil && invoices != nil {
			mu.Lock()
			results.Invoices = invoices
			results.TotalResults += len(invoices)
			mu.Unlock()
		}
	}()

	// Wait for all searches to complete
	wg.Wait()

	return results, nil
}
