package repository

import (
	"context"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
)

// StatsRepository defines the interface for statistics data operations
type StatsRepository interface {
	// GetDashboardStats retrieves all dashboard statistics
	GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error)
}
