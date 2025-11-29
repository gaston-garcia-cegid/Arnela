package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type employeeRepository struct {
	db *sqlx.DB
}

// NewEmployeeRepository creates a new instance of EmployeeRepository
func NewEmployeeRepository(db *sqlx.DB) repository.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	query := `
		INSERT INTO employees (
			id, user_id, first_name, last_name, email, phone, dni,
			date_of_birth, position, specialties, is_active, hire_date,
			notes, avatar_color
		) VALUES (
			:id, :user_id, :first_name, :last_name, :email, :phone, :dni,
			:date_of_birth, :position, :specialties, :is_active, :hire_date,
			:notes, :avatar_color
		)
		RETURNING created_at, updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, employee)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				if pqErr.Constraint == "employees_email_key" {
					return repository.ErrEmailAlreadyExists
				}
				if pqErr.Constraint == "employees_dni_key" {
					return repository.ErrDNIAlreadyExists
				}
			}
		}
		return fmt.Errorf("failed to create employee: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&employee.CreatedAt, &employee.UpdatedAt)
	}

	return errors.New("failed to retrieve created_at and updated_at")
}

func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT id, user_id, first_name, last_name, email, phone, dni,
		       date_of_birth, position, specialties, is_active, hire_date,
		       notes, avatar_color, created_at, updated_at, deleted_at
		FROM employees
		WHERE id = $1 AND deleted_at IS NULL
	`

	var employee domain.Employee
	err := r.db.GetContext(ctx, &employee, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee by ID: %w", err)
	}

	return &employee, nil
}

func (r *employeeRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Employee, error) {
	query := `
		SELECT id, user_id, first_name, last_name, email, phone, dni,
		       date_of_birth, position, specialties, is_active, hire_date,
		       notes, avatar_color, created_at, updated_at, deleted_at
		FROM employees
		WHERE user_id = $1 AND deleted_at IS NULL
	`

	var employee domain.Employee
	err := r.db.GetContext(ctx, &employee, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee by user ID: %w", err)
	}

	return &employee, nil
}

func (r *employeeRepository) GetByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	query := `
		SELECT id, user_id, first_name, last_name, email, phone, dni,
		       date_of_birth, position, specialties, is_active, hire_date,
		       notes, avatar_color, created_at, updated_at, deleted_at
		FROM employees
		WHERE email = $1 AND deleted_at IS NULL
	`

	var employee domain.Employee
	err := r.db.GetContext(ctx, &employee, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee by email: %w", err)
	}

	return &employee, nil
}

func (r *employeeRepository) GetByDNI(ctx context.Context, dni string) (*domain.Employee, error) {
	query := `
		SELECT id, user_id, first_name, last_name, email, phone, dni,
		       date_of_birth, position, specialties, is_active, hire_date,
		       notes, avatar_color, created_at, updated_at, deleted_at
		FROM employees
		WHERE dni = $1 AND deleted_at IS NULL
	`

	var employee domain.Employee
	err := r.db.GetContext(ctx, &employee, query, dni)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("failed to get employee by DNI: %w", err)
	}

	return &employee, nil
}

func (r *employeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	query := `
		UPDATE employees
		SET first_name = :first_name,
		    last_name = :last_name,
		    email = :email,
		    phone = :phone,
		    dni = :dni,
		    date_of_birth = :date_of_birth,
		    position = :position,
		    specialties = :specialties,
		    is_active = :is_active,
		    hire_date = :hire_date,
		    notes = :notes,
		    avatar_color = :avatar_color
		WHERE id = :id AND deleted_at IS NULL
		RETURNING updated_at
	`

	rows, err := r.db.NamedQueryContext(ctx, query, employee)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				if pqErr.Constraint == "employees_email_key" {
					return repository.ErrEmailAlreadyExists
				}
				if pqErr.Constraint == "employees_dni_key" {
					return repository.ErrDNIAlreadyExists
				}
			}
		}
		return fmt.Errorf("failed to update employee: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return repository.ErrEmployeeNotFound
	}

	return rows.Scan(&employee.UpdatedAt)
}

func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE employees
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return repository.ErrEmployeeNotFound
	}

	return nil
}

func (r *employeeRepository) List(ctx context.Context, limit, offset int) ([]*domain.Employee, error) {
	query := `
		SELECT id, user_id, first_name, last_name, email, phone, dni,
		       date_of_birth, position, specialties, is_active, hire_date,
		       notes, avatar_color, created_at, updated_at, deleted_at
		FROM employees
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var employees []*domain.Employee
	err := r.db.SelectContext(ctx, &employees, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list employees: %w", err)
	}

	return employees, nil
}

func (r *employeeRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM employees WHERE deleted_at IS NULL`

	var count int
	err := r.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return count, nil
}

func (r *employeeRepository) GetBySpecialty(ctx context.Context, specialty string) ([]*domain.Employee, error) {
	query := `
		SELECT id, user_id, first_name, last_name, email, phone, dni,
		       date_of_birth, position, specialties, is_active, hire_date,
		       notes, avatar_color, created_at, updated_at, deleted_at
		FROM employees
		WHERE $1 = ANY(specialties) AND deleted_at IS NULL AND is_active = true
		ORDER BY first_name, last_name
	`

	var employees []*domain.Employee
	err := r.db.SelectContext(ctx, &employees, query, specialty)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by specialty: %w", err)
	}

	return employees, nil
}

func (r *employeeRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM employees WHERE email = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists, nil
}

func (r *employeeRepository) DNIExists(ctx context.Context, dni string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM employees WHERE dni = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.GetContext(ctx, &exists, query, dni)
	if err != nil {
		return false, fmt.Errorf("failed to check if DNI exists: %w", err)
	}

	return exists, nil
}
