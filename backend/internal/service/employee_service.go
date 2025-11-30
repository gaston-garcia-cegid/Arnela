package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// EmployeeService defines the interface for employee business logic
type EmployeeService interface {
	CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*domain.Employee, error)
	GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error)
	GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error)
	UpdateEmployee(ctx context.Context, id uuid.UUID, req UpdateEmployeeRequest) (*domain.Employee, error)
	DeleteEmployee(ctx context.Context, id uuid.UUID) error
	ListEmployees(ctx context.Context, limit, offset int) ([]*domain.Employee, int, error)
	GetEmployeesBySpecialty(ctx context.Context, specialty string) ([]*domain.Employee, error)
}

type employeeService struct {
	repo     repository.EmployeeRepository
	userRepo repository.UserRepository
}

// NewEmployeeService creates a new instance of EmployeeService
func NewEmployeeService(repo repository.EmployeeRepository, userRepo repository.UserRepository) EmployeeService {
	return &employeeService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// CreateEmployeeRequest represents the request to create a new employee
type CreateEmployeeRequest struct {
	UserID      *uuid.UUID `json:"userId"`
	FirstName   string     `json:"firstName" binding:"required"`
	LastName    string     `json:"lastName" binding:"required"`
	Email       string     `json:"email" binding:"required,email"`
	Phone       string     `json:"phone" binding:"required"`
	DNI         string     `json:"dni" binding:"required"`
	Specialty   string     `json:"specialty" binding:"required"` // Single specialty, converted to array internally
	HireDate    string     `json:"hireDate" binding:"required"`  // ISO 8601 date string
	Notes       string     `json:"notes"`                        // Optional notes
	AvatarColor string     `json:"avatarColor"`
}

// UpdateEmployeeRequest represents the request to update an employee
type UpdateEmployeeRequest struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	DNI         string `json:"dni"`
	Specialty   string `json:"specialty"` // Single specialty
	HireDate    string `json:"hireDate"`  // ISO 8601 date string
	Notes       string `json:"notes"`     // Optional notes
	IsActive    *bool  `json:"isActive"`
	AvatarColor string `json:"avatarColor"`
}

var (
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrInvalidPhone     = errors.New("invalid phone format (Spanish format required)")
	ErrInvalidDNI       = errors.New("invalid DNI format")
	ErrEmailInUse       = errors.New("email is already in use")
	ErrDNIInUse         = errors.New("DNI is already in use")
	ErrEmployeeNotFound = errors.New("employee not found")
)

func (s *employeeService) CreateEmployee(ctx context.Context, req CreateEmployeeRequest) (*domain.Employee, error) {
	// Validate formats first
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := validateSpanishPhone(req.Phone); err != nil {
		return nil, err
	}

	if err := validateSpanishDNI(req.DNI); err != nil {
		return nil, err
	}

	// Check for duplicates
	emailExists, err := s.repo.EmailExists(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	}
	if emailExists {
		return nil, ErrEmailInUse
	}

	dniExists, err := s.repo.DNIExists(ctx, req.DNI)
	if err != nil {
		return nil, fmt.Errorf("failed to check DNI existence: %w", err)
	}
	if dniExists {
		return nil, ErrDNIInUse
	}

	// Parse hireDate
	var hireDate *time.Time
	if req.HireDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.HireDate)
		if err != nil {
			return nil, fmt.Errorf("invalid hire date format, expected YYYY-MM-DD: %w", err)
		}
		hireDate = &parsedDate
	}

	// Convert single specialty to array
	var specialties domain.StringArray
	if req.Specialty != "" {
		specialties = domain.StringArray{strings.TrimSpace(req.Specialty)}
	}

	// Create user for employee (email as username, DNI as password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.DNI), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(hashedPassword),
		FirstName:    strings.TrimSpace(req.FirstName),
		LastName:     strings.TrimSpace(req.LastName),
		Role:         domain.RoleEmployee,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Create employee entity linked to user
	employee := &domain.Employee{
		ID:          uuid.New(),
		UserID:      &user.ID,
		FirstName:   strings.TrimSpace(req.FirstName),
		LastName:    strings.TrimSpace(req.LastName),
		Email:       strings.ToLower(strings.TrimSpace(req.Email)),
		Phone:       strings.TrimSpace(req.Phone),
		DNI:         strings.ToUpper(strings.TrimSpace(req.DNI)),
		Position:    strings.TrimSpace(req.Specialty), // Use specialty as position
		Specialties: specialties,
		HireDate:    hireDate,
		Notes:       strings.TrimSpace(req.Notes),
		IsActive:    true,
		AvatarColor: req.AvatarColor,
	}

	if employee.AvatarColor == "" {
		employee.AvatarColor = "#6366F1" // Default color
	}

	// Save to repository
	if err := s.repo.Create(ctx, employee); err != nil {
		// Rollback: delete user if employee creation fails
		if deleteErr := s.userRepo.Delete(ctx, user.ID); deleteErr != nil {
			fmt.Printf("[ERROR] Failed to rollback user creation: %v\n", deleteErr)
		}
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, ErrEmailInUse
		}
		if errors.Is(err, repository.ErrDNIAlreadyExists) {
			return nil, ErrDNIInUse
		}
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	return employee, nil
}

func (s *employeeService) GetEmployee(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}
	return employee, nil
}

func (s *employeeService) GetEmployeeByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error) {
	employee, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee by user ID: %w", err)
	}
	return employee, nil
}

func (s *employeeService) UpdateEmployee(ctx context.Context, id uuid.UUID, req UpdateEmployeeRequest) (*domain.Employee, error) {
	// Get existing employee
	employee, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Update fields if provided
	if req.FirstName != "" {
		employee.FirstName = strings.TrimSpace(req.FirstName)
	}
	if req.LastName != "" {
		employee.LastName = strings.TrimSpace(req.LastName)
	}
	if req.Email != "" {
		if err := validateEmail(req.Email); err != nil {
			return nil, err
		}
		emailExists, err := s.repo.EmailExists(ctx, req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		}
		if emailExists {
			existing, _ := s.repo.GetByEmail(ctx, req.Email)
			if existing != nil && existing.ID != id {
				return nil, ErrEmailInUse
			}
		}
		employee.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}
	if req.Phone != "" {
		if err := validateSpanishPhone(req.Phone); err != nil {
			return nil, err
		}
		employee.Phone = strings.TrimSpace(req.Phone)
	}
	if req.DNI != "" {
		if err := validateSpanishDNI(req.DNI); err != nil {
			return nil, err
		}
		dniExists, err := s.repo.DNIExists(ctx, req.DNI)
		if err != nil {
			return nil, fmt.Errorf("failed to check DNI existence: %w", err)
		}
		if dniExists {
			existing, _ := s.repo.GetByDNI(ctx, req.DNI)
			if existing != nil && existing.ID != id {
				return nil, ErrDNIInUse
			}
		}
		employee.DNI = strings.ToUpper(strings.TrimSpace(req.DNI))
	}
	if req.Specialty != "" {
		employee.Position = strings.TrimSpace(req.Specialty)
		employee.Specialties = domain.StringArray{strings.TrimSpace(req.Specialty)}
	}
	if req.HireDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.HireDate)
		if err != nil {
			return nil, fmt.Errorf("invalid hire date format, expected YYYY-MM-DD: %w", err)
		}
		employee.HireDate = &parsedDate
	}
	if req.Notes != "" {
		employee.Notes = strings.TrimSpace(req.Notes)
	}
	if req.IsActive != nil {
		employee.IsActive = *req.IsActive
	}
	if req.AvatarColor != "" {
		employee.AvatarColor = req.AvatarColor
	}

	// Save updates
	if err := s.repo.Update(ctx, employee); err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return nil, ErrEmployeeNotFound
		}
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, ErrEmailInUse
		}
		if errors.Is(err, repository.ErrDNIAlreadyExists) {
			return nil, ErrDNIInUse
		}
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	return employee, nil
}

func (s *employeeService) DeleteEmployee(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			return ErrEmployeeNotFound
		}
		return fmt.Errorf("failed to delete employee: %w", err)
	}
	return nil
}

func (s *employeeService) ListEmployees(ctx context.Context, limit, offset int) ([]*domain.Employee, int, error) {
	employees, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}

	count, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return employees, count, nil
}

func (s *employeeService) GetEmployeesBySpecialty(ctx context.Context, specialty string) ([]*domain.Employee, error) {
	employees, err := s.repo.GetBySpecialty(ctx, specialty)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by specialty: %w", err)
	}
	return employees, nil
}

// Validation helpers
func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return ErrInvalidEmail
	}
	return nil
}

func validateSpanishPhone(phone string) error {
	// Spanish phone: +34 XXX XXX XXX or 6XX XXX XXX or 9XX XXX XXX
	phoneRegex := regexp.MustCompile(`^(\+34|0034)?[6-9]\d{8}$`)
	cleanPhone := regexp.MustCompile(`\s+`).ReplaceAllString(phone, "")
	if !phoneRegex.MatchString(cleanPhone) {
		return ErrInvalidPhone
	}
	return nil
}

func validateSpanishDNI(dni string) error {
	// Spanish DNI: 8 digits + 1 letter
	dniRegex := regexp.MustCompile(`^\d{8}[A-Z]$`)
	cleanDNI := strings.ToUpper(strings.TrimSpace(dni))
	if !dniRegex.MatchString(cleanDNI) {
		return ErrInvalidDNI
	}

	// Validate letter
	letters := "TRWAGMYFPDXBNJZSQVHLCKE"
	digits := cleanDNI[:8]
	letter := cleanDNI[8:]

	var num int
	fmt.Sscanf(digits, "%d", &num)
	expectedLetter := string(letters[num%23])

	if letter != expectedLetter {
		return ErrInvalidDNI
	}

	return nil
}
