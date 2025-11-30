package service

import (
	"context"
	"fmt"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/repository"
)

// StatsService handles business logic for statistics
type StatsService struct {
	statsRepo repository.StatsRepository
}

// NewStatsService creates a new StatsService instance
func NewStatsService(statsRepo repository.StatsRepository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

// GetDashboardStats retrieves all dashboard statistics
// Returns aggregated statistics for clients, employees, and appointments
func (s *StatsService) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	stats, err := s.statsRepo.GetDashboardStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard stats: %w", err)
	}

	return stats, nil
}
