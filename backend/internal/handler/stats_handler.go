package handler

import (
	"net/http"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

// StatsHandler handles statistics-related endpoints
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetDashboardStats retrieves all dashboard statistics
// @Summary      Get dashboard statistics
// @Description  Retrieves aggregated statistics for clients, employees, and appointments (admin/employee only)
// @Tags         stats
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} domain.DashboardStats
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /api/v1/stats/dashboard [get]
func (h *StatsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.statsService.GetDashboardStats(c.Request.Context())
	if err != nil {
		appErr := pkgerrors.NewInternalError("Error al obtener estadísticas")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	c.JSON(http.StatusOK, stats)
}
