package postgres

import (
	"context"
	"fmt"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type statsRepository struct {
	db *sqlx.DB
}

// NewStatsRepository creates a new StatsRepository instance
func NewStatsRepository(db *sqlx.DB) repository.StatsRepository {
	return &statsRepository{db: db}
}

// GetDashboardStats retrieves all dashboard statistics in a single optimized query
func (r *statsRepository) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	query := `
		SELECT
			-- Client stats
			COUNT(DISTINCT c.id) FILTER (WHERE c.deleted_at IS NULL) as total_clients,
			COUNT(DISTINCT c.id) FILTER (WHERE c.deleted_at IS NULL AND c.is_active = true) as active_clients,
			COUNT(DISTINCT c.id) FILTER (WHERE c.deleted_at IS NULL AND c.is_active = false) as inactive_clients,
			
			-- Employee stats
			COUNT(DISTINCT e.id) FILTER (WHERE e.deleted_at IS NULL) as total_employees,
			COUNT(DISTINCT e.id) FILTER (WHERE e.deleted_at IS NULL AND e.is_active = true) as active_employees,
			COUNT(DISTINCT e.id) FILTER (WHERE e.deleted_at IS NULL AND e.is_active = false) as inactive_employees,
			
			-- Appointment stats
			COUNT(DISTINCT a.id) FILTER (WHERE a.deleted_at IS NULL) as total_appointments,
			COUNT(DISTINCT a.id) FILTER (WHERE a.deleted_at IS NULL AND a.status = 'pending') as pending_appointments,
			COUNT(DISTINCT a.id) FILTER (WHERE a.deleted_at IS NULL AND a.status = 'confirmed') as confirmed_appointments,
			COUNT(DISTINCT a.id) FILTER (WHERE a.deleted_at IS NULL AND a.status = 'completed') as completed_appointments,
			COUNT(DISTINCT a.id) FILTER (WHERE a.deleted_at IS NULL AND a.status = 'cancelled') as cancelled_appointments
		FROM
			clients c
		CROSS JOIN
			employees e
		CROSS JOIN
			appointments a
	`

	var result struct {
		TotalClients    int `db:"total_clients"`
		ActiveClients   int `db:"active_clients"`
		InactiveClients int `db:"inactive_clients"`

		TotalEmployees    int `db:"total_employees"`
		ActiveEmployees   int `db:"active_employees"`
		InactiveEmployees int `db:"inactive_employees"`

		TotalAppointments     int `db:"total_appointments"`
		PendingAppointments   int `db:"pending_appointments"`
		ConfirmedAppointments int `db:"confirmed_appointments"`
		CompletedAppointments int `db:"completed_appointments"`
		CancelledAppointments int `db:"cancelled_appointments"`
	}

	err := r.db.GetContext(ctx, &result, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	stats := &domain.DashboardStats{
		Clients: domain.ClientStats{
			Total:    result.TotalClients,
			Active:   result.ActiveClients,
			Inactive: result.InactiveClients,
		},
		Employees: domain.EmployeeStats{
			Total:    result.TotalEmployees,
			Active:   result.ActiveEmployees,
			Inactive: result.InactiveEmployees,
		},
		Appointments: domain.AppointmentStats{
			Total:     result.TotalAppointments,
			Pending:   result.PendingAppointments,
			Confirmed: result.ConfirmedAppointments,
			Completed: result.CompletedAppointments,
			Cancelled: result.CancelledAppointments,
		},
	}

	return stats, nil
}
