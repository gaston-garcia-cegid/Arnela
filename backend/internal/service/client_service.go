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
	userRepo   repository.UserRepository // ✅ AÑADIDO
}

func NewClientService(clientRepo repository.ClientRepository, userRepo repository.UserRepository) ClientService {
	return &clientService{
		clientRepo: clientRepo,
		userRepo:   userRepo,
	}
}

func (s *clientService) CreateClient(ctx context.Context, req CreateClientRequest) (*domain.Client, error) {
	log.Printf("[DEBUG] Creating client with email: %s", req.Email)

	// ✅ PASO 1: Verificar que el email no exista en clients
	if exists, err := s.clientRepo.EmailExists(ctx, req.Email, req.UserID); err != nil {
		return nil, fmt.Errorf("failed to check email existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("email already registered")
	}

	// ✅ PASO 2: Verificar que el email no exista en users
	if exists, err := s.userRepo.EmailExists(ctx, req.Email); err != nil {
		return nil, fmt.Errorf("failed to check user email existence: %w", err)
	} else if exists {
		return nil, fmt.Errorf("email already registered as user")
	}

	// ✅ PASO 3: Verificar DNI único
	if req.DNI != "" {
		if exists, err := s.clientRepo.DNIExists(ctx, req.DNI, req.UserID); err != nil {
			return nil, fmt.Errorf("failed to check DNI existence: %w", err)
		} else if exists {
			return nil, fmt.Errorf("DNI already registered")
		}
	}

	// ✅ PASO 4: Verificar NIF único (si se proporciona)
	if req.NIF != "" {
		if exists, err := s.clientRepo.NIFExists(ctx, req.NIF, req.UserID); err != nil {
			return nil, fmt.Errorf("failed to check NIF existence: %w", err)
		} else if exists {
			return nil, fmt.Errorf("NIF already registered")
		}
	}

	// ✅ PASO 5: Crear usuario con DNI como contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.DNI), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         domain.RoleClient, // ✅ Rol de cliente
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	log.Printf("[DEBUG] User created: ID=%s, Email=%s, Role=%s", user.ID, user.Email, user.Role)

	// ✅ PASO 6: Crear cliente vinculado al usuario
	client := &domain.Client{
		ID:        uuid.New(),
		UserID:    user.ID, // ✅ Vincular con usuario
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		DNI:       req.DNI,
		NIF:       req.NIF,
		Address:   domain.Address{Street: req.Address},
		Notes:     req.Notes,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.clientRepo.Create(ctx, client); err != nil {
		// ✅ ROLLBACK: Si falla la creación del cliente, eliminar usuario
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
		if exists, err := s.clientRepo.EmailExists(ctx, *req.Email, req.UserID); err != nil {
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
	if req.DNI != nil {
		if exists, err := s.clientRepo.DNIExists(ctx, *req.DNI, req.UserID); err != nil {
			return nil, fmt.Errorf("failed to check DNI existence: %w", err)
		} else if exists && *req.DNI != client.DNI {
			return nil, fmt.Errorf("DNI already registered")
		}
		client.DNI = *req.DNI
	}
	if req.NIF != nil {
		if exists, err := s.clientRepo.NIFExists(ctx, *req.NIF, req.UserID); err != nil {
			return nil, fmt.Errorf("failed to check NIF existence: %w", err)
		} else if exists && *req.NIF != client.NIF {
			return nil, fmt.Errorf("NIF already registered")
		}
		client.NIF = *req.NIF
	}
	if req.Address != nil {
		client.Address = domain.Address{Street: *req.Address}
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
	// Get paginated clients from repository
	clients, err := s.clientRepo.List(ctx, filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	// Calculate total and pagination metadata
	total := len(clients)
	totalPages := (total + pageSize - 1) / pageSize

	return &ClientListResponse{
		Clients:    clients,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
