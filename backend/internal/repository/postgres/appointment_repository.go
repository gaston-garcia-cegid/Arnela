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

type appointmentRepository struct {
	db *sqlx.DB
}

func NewAppointmentRepository(db *sqlx.DB) repository.AppointmentRepository {
	return &appointmentRepository{db: db}
}

// ✅ Columnas base para SELECT
const appointmentColumns = `
    id, client_id, therapist_id, title, description,
    start_time, end_time, duration_minutes, status,
    notes, cancellation_reason, google_calendar_event_id,
    created_by, created_at, updated_at, deleted_at
`

func (r *appointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	query := `
        INSERT INTO appointments (
            id, client_id, therapist_id, title, description,
            start_time, end_time, duration_minutes, status,
            notes, created_by, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	_, err := r.db.ExecContext(ctx, query,
		appointment.ID,
		appointment.ClientID,
		appointment.TherapistID,
		appointment.Title,
		appointment.Description,
		appointment.StartTime,
		appointment.EndTime,
		appointment.DurationMinutes,
		appointment.Status,
		appointment.Notes,
		appointment.CreatedBy,
		appointment.CreatedAt,
		appointment.UpdatedAt,
	)

	return err
}

func (r *appointmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	var appointment domain.Appointment
	query := fmt.Sprintf(`
        SELECT %s
        FROM appointments
        WHERE id = $1 AND deleted_at IS NULL
    `, appointmentColumns)

	err := r.db.GetContext(ctx, &appointment, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("appointment not found")
		}
		return nil, fmt.Errorf("failed to get appointment: %w", err)
	}

	return &appointment, nil
}

func (r *appointmentRepository) GetByIDWithRelations(ctx context.Context, id uuid.UUID) (*domain.Appointment, error) {
	appointment, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load therapist data (mock)
	therapist := domain.GetTherapistByID(appointment.TherapistID)
	appointment.Therapist = therapist

	// Load client data
	var client domain.Client
	clientQuery := `
        SELECT id, user_id, email, first_name, last_name, phone, dni
        FROM clients
        WHERE id = $1 AND deleted_at IS NULL
    `
	err = r.db.GetContext(ctx, &client, clientQuery, appointment.ClientID)
	if err == nil {
		appointment.Client = &client
	}

	return appointment, nil
}

func (r *appointmentRepository) Update(ctx context.Context, appointment *domain.Appointment) error {
	query := `
		UPDATE appointments SET
			therapist_id = $1,
			title = $2,
			description = $3,
			start_time = $4,
			end_time = $5,
			duration_minutes = $6,
			status = $7,
			notes = $8,
			cancellation_reason = $9,
			google_calendar_event_id = $10,
			updated_at = $11
		WHERE id = $12 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query,
		appointment.TherapistID,
		appointment.Title,
		appointment.Description,
		appointment.StartTime,
		appointment.EndTime,
		appointment.DurationMinutes,
		appointment.Status,
		appointment.Notes,
		appointment.CancellationReason,
		appointment.GoogleCalendarEventID,
		time.Now(),
		appointment.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update appointment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("appointment not found")
	}

	return nil
}

func (r *appointmentRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AppointmentStatus) error {
	query := `
		UPDATE appointments SET
			status = $1,
			updated_at = $2
		WHERE id = $3 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update appointment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("appointment not found")
	}

	return nil
}

func (r *appointmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
        UPDATE appointments SET deleted_at = $1
        WHERE id = $2 AND deleted_at IS NULL
    `

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to delete appointment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("appointment not found")
	}

	return nil
}

func (r *appointmentRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time, therapistID *string) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment

	query := fmt.Sprintf(`
        SELECT %s
        FROM appointments
        WHERE deleted_at IS NULL
        AND start_time >= $1
        AND end_time <= $2
    `, appointmentColumns)

	args := []interface{}{startDate, endDate}

	if therapistID != nil {
		query += " AND therapist_id = $3"
		args = append(args, *therapistID)
	}

	query += " ORDER BY start_time ASC"

	err := r.db.SelectContext(ctx, &appointments, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by date range: %w", err)
	}

	return appointments, nil
}

func (r *appointmentRepository) List(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment

	query := fmt.Sprintf(`
		SELECT %s
		FROM appointments
		WHERE deleted_at IS NULL
	`, appointmentColumns)

	args := []interface{}{}
	argCount := 0

	conditions := []string{}

	if filters.ClientID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("client_id = $%d", argCount))
		args = append(args, *filters.ClientID)
	}

	if filters.TherapistID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("therapist_id = $%d", argCount))
		args = append(args, *filters.TherapistID)
	}

	if filters.Status != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *filters.Status)
	}

	if filters.StartDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", argCount))
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("end_time <= $%d", argCount))
		args = append(args, *filters.EndDate)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY start_time DESC"

	if filters.PageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount+1, argCount+2)
		args = append(args, filters.PageSize, (filters.Page-1)*filters.PageSize)
	}

	err := r.db.SelectContext(ctx, &appointments, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list appointments: %w", err)
	}

	return appointments, nil
}

func (r *appointmentRepository) ListWithRelations(ctx context.Context, filters domain.AppointmentFilter) ([]*domain.Appointment, error) {
	appointments, err := r.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	// Load relations for each appointment
	for _, apt := range appointments {
		therapist := domain.GetTherapistByID(apt.TherapistID)
		apt.Therapist = therapist

		var client domain.Client
		clientQuery := `
			SELECT id, user_id, email, first_name, last_name, phone, dni
			FROM clients
			WHERE id = $1 AND deleted_at IS NULL
		`
		err = r.db.GetContext(ctx, &client, clientQuery, apt.ClientID)
		if err == nil {
			apt.Client = &client
		}
	}

	return appointments, nil
}

func (r *appointmentRepository) Count(ctx context.Context, filters domain.AppointmentFilter) (int, error) {
	query := `SELECT COUNT(*) FROM appointments WHERE deleted_at IS NULL`
	args := []interface{}{}
	argCount := 0

	conditions := []string{}

	if filters.ClientID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("client_id = $%d", argCount))
		args = append(args, *filters.ClientID)
	}

	if filters.TherapistID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("therapist_id = $%d", argCount))
		args = append(args, *filters.TherapistID)
	}

	if filters.Status != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *filters.Status)
	}

	if filters.StartDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", argCount))
		args = append(args, *filters.StartDate)
	}

	if filters.EndDate != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("end_time <= $%d", argCount))
		args = append(args, *filters.EndDate)
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to count appointments: %w", err)
	}

	return count, nil
}

func (r *appointmentRepository) CheckOverlap(ctx context.Context, therapistID string, startTime, endTime time.Time, excludeID *uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM appointments
		WHERE therapist_id = $1
		AND deleted_at IS NULL
		AND status != 'cancelled'
		AND (
			(start_time < $3 AND end_time > $2)
			OR (start_time >= $2 AND start_time < $3)
			OR (end_time > $2 AND end_time <= $3)
		)
	`

	args := []interface{}{therapistID, startTime, endTime}

	if excludeID != nil {
		query += " AND id != $4"
		args = append(args, *excludeID)
	}

	var count int
	err := r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return false, fmt.Errorf("failed to check appointment overlap: %w", err)
	}

	return count > 0, nil
}

func (r *appointmentRepository) GetByClientID(ctx context.Context, clientID uuid.UUID, page, pageSize int) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment

	query := fmt.Sprintf(`
		SELECT %s
		FROM appointments
		WHERE client_id = $1 AND deleted_at IS NULL
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`, appointmentColumns)

	offset := (page - 1) * pageSize
	err := r.db.SelectContext(ctx, &appointments, query, clientID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by client ID: %w", err)
	}

	return appointments, nil
}

func (r *appointmentRepository) GetByTherapistID(ctx context.Context, therapistID string, page, pageSize int) ([]*domain.Appointment, error) {
	var appointments []*domain.Appointment

	query := fmt.Sprintf(`
		SELECT %s
		FROM appointments
		WHERE therapist_id = $1 AND deleted_at IS NULL
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`, appointmentColumns)

	offset := (page - 1) * pageSize
	err := r.db.SelectContext(ctx, &appointments, query, therapistID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments by therapist ID: %w", err)
	}

	return appointments, nil
}
