package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
)

// ClientService handles client business logic
type ClientService struct {
	clientRepo repository.ClientRepository
}

// NewClientService creates a new ClientService
func NewClientService(clientRepo repository.ClientRepository) *ClientService {
	return &ClientService{
		clientRepo: clientRepo,
	}
}

// CreateClient creates a new client with validations
func (s *ClientService) CreateClient(ctx context.Context, req CreateClientRequest) (*domain.Client, error) {
	// Validate email format
	if !isValidEmail(req.Email) {
		return nil, fmt.Errorf("invalid email format")
	}

	// Validate phone format (Spanish phone)
	if !isValidSpanishPhone(req.Phone) {
		return nil, fmt.Errorf("invalid phone format (expected Spanish phone number)")
	}

	// Validate NIF format (required)
	if !isValidSpanishDNI(req.NIF) {
		return nil, fmt.Errorf("invalid NIF format")
	}

	// Validate DNI format if provided (optional)
	if req.DNI != "" && !isValidSpanishDNI(req.DNI) {
		return nil, fmt.Errorf("invalid DNI/NIE format")
	}

	// Check if email already exists
	exists, err := s.clientRepo.EmailExists(ctx, req.Email, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// Check if NIF already exists
	exists, err = s.clientRepo.NIFExists(ctx, req.NIF, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to check NIF existence: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("NIF already registered")
	}

	// Check if DNI already exists (if provided)
	if req.DNI != "" {
		exists, err = s.clientRepo.DNIExists(ctx, req.DNI, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to check DNI existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("DNI already registered")
		}
	}

	// Parse date of birth if provided
	var dateOfBirth *time.Time
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format (expected YYYY-MM-DD)")
		}
		dateOfBirth = &dob
	}

	// Create client entity
	client := &domain.Client{
		ID:          uuid.New(),
		UserID:      req.UserID,
		FirstName:   strings.TrimSpace(req.FirstName),
		LastName:    strings.TrimSpace(req.LastName),
		Email:       strings.ToLower(strings.TrimSpace(req.Email)),
		Phone:       normalizePhone(req.Phone),
		NIF:         strings.ToUpper(strings.TrimSpace(req.NIF)),
		DNI:         strings.ToUpper(strings.TrimSpace(req.DNI)),
		DateOfBirth: dateOfBirth,
		Address:     strings.TrimSpace(req.Address),
		City:        strings.TrimSpace(req.City),
		PostalCode:  strings.TrimSpace(req.PostalCode),
		Province:    strings.TrimSpace(req.Province),
		IsActive:    true,
		Notes:       strings.TrimSpace(req.Notes),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to database
	if err := s.clientRepo.Create(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return client, nil
}

// GetClient retrieves a client by ID
func (s *ClientService) GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	client, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	return client, nil
}

// GetClientByUserID retrieves a client by their associated user ID
func (s *ClientService) GetClientByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error) {
	client, err := s.clientRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	return client, nil
}

// UpdateClient updates a client's information
func (s *ClientService) UpdateClient(ctx context.Context, id uuid.UUID, req UpdateClientRequest) (*domain.Client, error) {
	// Get existing client
	client, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	// Update fields if provided
	if req.FirstName != nil {
		client.FirstName = strings.TrimSpace(*req.FirstName)
	}

	if req.LastName != nil {
		client.LastName = strings.TrimSpace(*req.LastName)
	}

	if req.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*req.Email))
		if !isValidEmail(email) {
			return nil, fmt.Errorf("invalid email format")
		}

		// Check if email already exists (excluding current client)
		exists, err := s.clientRepo.EmailExists(ctx, email, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("email already registered")
		}

		client.Email = email
	}

	if req.Phone != nil {
		phone := *req.Phone
		if !isValidSpanishPhone(phone) {
			return nil, fmt.Errorf("invalid phone format")
		}
		client.Phone = normalizePhone(phone)
	}

	if req.NIF != nil {
		nif := strings.ToUpper(strings.TrimSpace(*req.NIF))
		if !isValidSpanishDNI(nif) {
			return nil, fmt.Errorf("invalid NIF format")
		}

		// Check if NIF already exists (excluding current client)
		exists, err := s.clientRepo.NIFExists(ctx, nif, &id)
		if err != nil {
			return nil, fmt.Errorf("failed to check NIF existence: %w", err)
		}
		if exists {
			return nil, fmt.Errorf("NIF already registered")
		}

		client.NIF = nif
	}

	if req.DNI != nil {
		dni := strings.ToUpper(strings.TrimSpace(*req.DNI))
		if dni != "" && !isValidSpanishDNI(dni) {
			return nil, fmt.Errorf("invalid DNI/NIE format")
		}

		// Check if DNI already exists (excluding current client)
		if dni != "" {
			exists, err := s.clientRepo.DNIExists(ctx, dni, &id)
			if err != nil {
				return nil, fmt.Errorf("failed to check DNI existence: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("DNI already registered")
			}
		}

		client.DNI = dni
	}

	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return nil, fmt.Errorf("invalid date of birth format (expected YYYY-MM-DD)")
		}
		client.DateOfBirth = &dob
	}

	if req.Address != nil {
		client.Address = strings.TrimSpace(*req.Address)
	}

	if req.City != nil {
		client.City = strings.TrimSpace(*req.City)
	}

	if req.PostalCode != nil {
		client.PostalCode = strings.TrimSpace(*req.PostalCode)
	}

	if req.Province != nil {
		client.Province = strings.TrimSpace(*req.Province)
	}

	if req.IsActive != nil {
		client.IsActive = *req.IsActive
	}

	if req.Notes != nil {
		client.Notes = strings.TrimSpace(*req.Notes)
	}

	client.UpdatedAt = time.Now()

	// Save changes
	if err := s.clientRepo.Update(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return client, nil
}

// DeleteClient soft-deletes a client
func (s *ClientService) DeleteClient(ctx context.Context, id uuid.UUID) error {
	// Verify client exists
	_, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	// Soft delete
	if err := s.clientRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}

// ListClients retrieves a paginated list of clients with filters
func (s *ClientService) ListClients(ctx context.Context, filters repository.ClientFilters, page, pageSize int) (*ClientListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	offset := (page - 1) * pageSize

	// Get total count
	total, err := s.clientRepo.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to count clients: %w", err)
	}

	// Get clients
	clients, err := s.clientRepo.List(ctx, filters, offset, pageSize)
	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}

	totalPages := (total + pageSize - 1) / pageSize

	return &ClientListResponse{
		Clients:    clients,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Validation helper functions

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidSpanishPhone(phone string) bool {
	// Remove spaces, dashes, and parentheses
	cleanPhone := regexp.MustCompile(`[\s\-\(\)]+`).ReplaceAllString(phone, "")

	// Spanish phone: starts with +34 or 34 or nothing, followed by 9 digits
	phoneRegex := regexp.MustCompile(`^(\+34|34)?[6789]\d{8}$`)
	return phoneRegex.MatchString(cleanPhone)
}

func normalizePhone(phone string) string {
	// Remove spaces, dashes, and parentheses
	cleanPhone := regexp.MustCompile(`[\s\-\(\)]+`).ReplaceAllString(phone, "")

	// Add +34 if not present
	if !strings.HasPrefix(cleanPhone, "+34") && !strings.HasPrefix(cleanPhone, "34") {
		cleanPhone = "+34" + cleanPhone
	} else if strings.HasPrefix(cleanPhone, "34") && !strings.HasPrefix(cleanPhone, "+34") {
		cleanPhone = "+" + cleanPhone
	}

	return cleanPhone
}

func isValidSpanishDNI(dni string) bool {
	// Remove spaces and convert to uppercase
	cleanDNI := strings.ToUpper(strings.ReplaceAll(dni, " ", ""))

	// DNI: 8 digits + 1 letter
	dniRegex := regexp.MustCompile(`^\d{8}[A-Z]$`)
	if dniRegex.MatchString(cleanDNI) {
		return validateDNILetter(cleanDNI)
	}

	// NIE: X/Y/Z + 7 digits + 1 letter
	nieRegex := regexp.MustCompile(`^[XYZ]\d{7}[A-Z]$`)
	if nieRegex.MatchString(cleanDNI) {
		return validateNIELetter(cleanDNI)
	}

	return false
}

func validateDNILetter(dni string) bool {
	letters := "TRWAGMYFPDXBNJZSQVHLCKE"
	number := dni[:8]
	letter := dni[8:]

	var num int
	fmt.Sscanf(number, "%d", &num)

	expectedLetter := string(letters[num%23])
	return letter == expectedLetter
}

func validateNIELetter(nie string) bool {
	letters := "TRWAGMYFPDXBNJZSQVHLCKE"

	// Replace first letter with corresponding number
	firstChar := nie[0]
	var prefix string
	switch firstChar {
	case 'X':
		prefix = "0"
	case 'Y':
		prefix = "1"
	case 'Z':
		prefix = "2"
	}

	number := prefix + nie[1:8]
	letter := nie[8:]

	var num int
	fmt.Sscanf(number, "%d", &num)

	expectedLetter := string(letters[num%23])
	return letter == expectedLetter
}
