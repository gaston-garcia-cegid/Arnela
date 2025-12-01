package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ClientServiceInterface defines the interface for client service
type ClientService interface {
	CreateClient(ctx context.Context, req CreateClientRequest) (*domain.Client, error)
	GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error)
	GetClientByEmail(ctx context.Context, email string) (*domain.Client, error)
	GetClientByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error)
	UpdateClient(ctx context.Context, id uuid.UUID, req UpdateClientRequest) (*domain.Client, error)
	DeleteClient(ctx context.Context, id uuid.UUID) error
	ListClients(ctx context.Context, filters repository.ClientFilters, page, pageSize int) (*ClientListResponse, error)
}

type clientService struct {
	clientRepo repository.ClientRepository
	userRepo   repository.UserRepository
}

func NewClientService(clientRepo repository.ClientRepository, userRepo repository.UserRepository) ClientServiceInterface {
	return &clientService{
		clientRepo: clientRepo,
		userRepo:   userRepo,
	}
}

func (s *clientService) CreateClient(ctx context.Context, req CreateClientRequest) (*domain.Client, error) {
	log.Printf("[DEBUG] Creating client with email: %s", req.Email)

	// Verificaciones de existencia
	if exists, err := s.clientRepo.EmailExists(ctx, req.Email, nil); err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("email already registered")
	}

	if exists, err := s.userRepo.EmailExists(ctx, req.Email); err != nil {
		return nil, fmt.Errorf("failed to check user email existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("email already registered as user")
	}

	if exists, err := s.clientRepo.DNICIFExists(ctx, req.DNICIF, nil); err != nil {
		return nil, fmt.Errorf("failed to check DNI/CIF existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("DNI/CIF already registered")
	}

	// Crear usuario con DNI/CIF como contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.DNICIF), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         domain.RoleClient,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("[DEBUG] User created: ID=%s, Email=%s, Role=%s", user.ID, user.Email, user.Role)

	// Crear cliente vinculado al usuario
	client := &domain.Client{
		ID:        uuid.New(),
		UserID:    user.ID,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		DNICIF:    req.DNICIF,
		Notes:     req.Notes,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// ✅ Usar helper para asignar address
	client.SetAddress(req.Address)

	if err := s.clientRepo.Create(ctx, client); err != nil {
		// Rollback: eliminar usuario
		if deleteErr := s.userRepo.Delete(ctx, user.ID); deleteErr != nil {
			log.Printf("[ERROR] Failed to rollback user creation: %v", deleteErr)
		}
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	log.Printf("[DEBUG] Client created: ID=%s, UserID=%s, Email=%s", client.ID, client.UserID, client.Email)

	return client, nil
}

func (s *clientService) GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	return s.clientRepo.GetByID(ctx, id)
}

func (s *clientService) GetClientByEmail(ctx context.Context, email string) (*domain.Client, error) {
	return s.clientRepo.GetByEmail(ctx, email)
}

func (s *clientService) GetClientByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error) {
	return s.clientRepo.GetByUserID(ctx, userID)
}

func (s *clientService) UpdateClient(ctx context.Context, id uuid.UUID, req UpdateClientRequest) (*domain.Client, error) {
	client, err := s.clientRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("client not found")
	}

	// Update fields if provided
	if req.Email != nil {
		if exists, err := s.clientRepo.EmailExists(ctx, *req.Email, &client.ID); err != nil {
			return nil, fmt.Errorf("failed to check email existence: %w", err)
		} else if exists && *req.Email != client.Email {
			return nil, fmt.Errorf("email already registered")
		}
		client.Email = *req.Email
	}
	if req.FirstName != nil {
		client.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		client.LastName = *req.LastName
	}
	if req.Phone != nil {
		client.Phone = *req.Phone
	}
	if req.DNICIF != nil {
		if exists, err := s.clientRepo.DNICIFExists(ctx, *req.DNICIF, &client.ID); err != nil {
			return nil, fmt.Errorf("failed to check DNI/CIF existence: %w", err)
		} else if exists && *req.DNICIF != client.DNICIF {
			return nil, fmt.Errorf("DNI/CIF already registered")
		}
		client.DNICIF = *req.DNICIF
	}
	if req.Address != nil {
		// ✅ Usar helper para actualizar address
		client.SetAddress(*req.Address)
	}
	if req.Notes != nil {
		client.Notes = *req.Notes
	}
	if req.IsActive != nil {
		client.IsActive = *req.IsActive
	}

	client.UpdatedAt = time.Now()

	if err := s.clientRepo.Update(ctx, client); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return client, nil
}

func (s *clientService) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return s.clientRepo.Delete(ctx, id)
}

func (s *clientService) ListClients(ctx context.Context, filters repository.ClientFilters, page, pageSize int) (*ClientListResponse, error) {
	// Get total count
	total, err := s.clientRepo.Count(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Get paginated clients
	clients, err := s.clientRepo.List(ctx, filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Calculate pagination metadata
	totalPages := (total + pageSize - 1) / pageSize

	return &ClientListResponse{
		Clients:    clients,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
