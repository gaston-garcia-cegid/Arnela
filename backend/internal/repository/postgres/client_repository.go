package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type clientRepository struct {
	db *sqlx.DB
}

func NewClientRepository(db *sqlx.DB) repository.ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) Create(ctx context.Context, client *domain.Client) error {
	query := `
		INSERT INTO clients (
			id, user_id, email, first_name, last_name, phone, dni, nif,
			address_street, address_city, address_province, address_postal_code, address_country,
			notes, is_active, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.db.ExecContext(ctx, query,
		client.ID,
		client.UserID, // ✅ AÑADIDO
		client.Email,
		client.FirstName,
		client.LastName,
		client.Phone,
		client.DNI,
		client.NIF,
		client.Address.Street,
		client.Address.City,
		client.Address.Province,
		client.Address.PostalCode,
		client.Address.Country,
		client.Notes,
		client.IsActive,
		client.CreatedAt,
		client.UpdatedAt,
	)
	return err
}

func (r *clientRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	var client domain.Client
	query := `
		SELECT id, user_id, email, first_name, last_name, phone, dni, nif,
			   address_street, address_city, address_province, address_postal_code, address_country,
			   notes, is_active, created_at, updated_at
		FROM clients
		WHERE id = $1 AND is_active = true
	`

	err := r.db.GetContext(ctx, &client, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return &client, nil
}

func (r *clientRepository) GetByEmail(ctx context.Context, email string) (*domain.Client, error) {
	var client domain.Client
	query := `
		SELECT id, user_id, email, first_name, last_name, phone, dni, nif,
			   address_street, address_city, address_province, address_postal_code, address_country,
			   notes, is_active, created_at, updated_at
		FROM clients
		WHERE email = $1 AND is_active = true
	`

	err := r.db.GetContext(ctx, &client, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client by email: %w", err)
	}

	return &client, nil
}

func (r *clientRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Client, error) {
	var client domain.Client
	query := `
		SELECT id, user_id, email, first_name, last_name, phone, dni, nif,
			   address_street, address_city, address_province, address_postal_code, address_country,
			   notes, is_active, created_at, updated_at
		FROM clients
		WHERE user_id = $1 AND is_active = true
	`

	err := r.db.GetContext(ctx, &client, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client by user ID: %w", err)
	}

	return &client, nil
}

func (r *clientRepository) GetByDNI(ctx context.Context, dni string) (*domain.Client, error) {
	var client domain.Client
	query := `
		SELECT id, user_id, email, first_name, last_name, phone, dni, nif,
			   address_street, address_city, address_province, address_postal_code, address_country,
			   notes, is_active, created_at, updated_at
		FROM clients
		WHERE dni = $1 AND is_active = true
	`

	err := r.db.GetContext(ctx, &client, query, dni)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client by DNI: %w", err)
	}

	return &client, nil
}

func (r *clientRepository) GetByNIF(ctx context.Context, nif string) (*domain.Client, error) {
	var client domain.Client
	query := `
		SELECT id, user_id, email, first_name, last_name, phone, dni, nif,
			   address_street, address_city, address_province, address_postal_code, address_country,
			   notes, is_active, created_at, updated_at
		FROM clients
		WHERE nif = $1 AND is_active = true
	`

	err := r.db.GetContext(ctx, &client, query, nif)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("client not found")
		}
		return nil, fmt.Errorf("failed to get client by NIF: %w", err)
	}

	return &client, nil
}

// Update updates a client's information
func (r *clientRepository) Update(ctx context.Context, client *domain.Client) error {
	query := `
		UPDATE clients SET
			user_id = :user_id,
			first_name = :first_name,
			last_name = :last_name,
			email = :email,
			phone = :phone,
			nif = :nif,
			dni = :dni,
			date_of_birth = :date_of_birth,
			address = :address,
			city = :city,
			postal_code = :postal_code,
			province = :province,
			is_active = :is_active,
			last_visit = :last_visit,
			notes = :notes,
			updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`

	result, err := r.db.NamedExecContext(ctx, query, client)
	if err != nil {
		return fmt.Errorf("failed to update client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("client not found")
	}

	return nil
}

// Delete soft-deletes a client
func (r *clientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE clients
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("client not found")
	}

	return nil
}

// List retrieves a paginated list of clients with optional filters
func (r *clientRepository) List(ctx context.Context, filters repository.ClientFilters, page, pageSize int) ([]*domain.Client, error) {
	var clients []*domain.Client

	query := `SELECT * FROM clients WHERE deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	// Build WHERE clause based on filters
	conditions := []string{}

	if filters.Search != "" {
		argCount++
		searchPattern := "%" + filters.Search + "%"
		conditions = append(conditions, fmt.Sprintf(`(
			first_name ILIKE $%d OR
			last_name ILIKE $%d OR
			email ILIKE $%d OR
			phone ILIKE $%d OR
			dni ILIKE $%d
		)`, argCount, argCount, argCount, argCount, argCount))
		args = append(args, searchPattern)
	}

	if filters.IsActive != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *filters.IsActive)
	}

	if filters.City != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("city ILIKE $%d", argCount))
		args = append(args, "%"+filters.City+"%")
	}

	if filters.Province != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("province ILIKE $%d", argCount))
		args = append(args, "%"+filters.Province+"%")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY last_name, first_name"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
	args = append(args, pageSize, (page-1)*pageSize)

	err := r.db.SelectContext(ctx, &clients, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list clients: %w", err)
	}

	return clients, nil
}

// Count returns the total number of clients matching the filters
func (r *clientRepository) Count(ctx context.Context, filters repository.ClientFilters) (int, error) {
	query := `SELECT COUNT(*) FROM clients WHERE deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	// Build WHERE clause based on filters (same as List)
	conditions := []string{}

	if filters.Search != "" {
		argCount++
		searchPattern := "%" + filters.Search + "%"
		conditions = append(conditions, fmt.Sprintf(`(
			first_name ILIKE $%d OR
			last_name ILIKE $%d OR
			email ILIKE $%d OR
			phone ILIKE $%d OR
			dni ILIKE $%d
		)`, argCount, argCount, argCount, argCount, argCount))
		args = append(args, searchPattern)
	}

	if filters.IsActive != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *filters.IsActive)
	}

	if filters.City != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("city ILIKE $%d", argCount))
		args = append(args, "%"+filters.City+"%")
	}

	if filters.Province != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("province ILIKE $%d", argCount))
		args = append(args, "%"+filters.Province+"%")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count clients: %w", err)
	}

	return count, nil
}

// EmailExists checks if an email is already registered
func (r *clientRepository) EmailExists(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM clients WHERE email = $1 AND deleted_at IS NULL`
	args := []interface{}{email}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	query += ")"

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}

	return exists, nil
}

// NIFExists checks if a NIF is already registered
func (r *clientRepository) NIFExists(ctx context.Context, nif string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM clients WHERE nif = $1 AND deleted_at IS NULL`
	args := []interface{}{nif}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	query += ")"

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check NIF existence: %w", err)
	}

	return exists, nil
}

// DNIExists checks if a DNI is already registered
func (r *clientRepository) DNIExists(ctx context.Context, dni string, excludeID *uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM clients WHERE dni = $1 AND dni != '' AND deleted_at IS NULL`
	args := []interface{}{dni}

	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}

	query += ")"

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check DNI existence: %w", err)
	}

	return exists, nil
}
