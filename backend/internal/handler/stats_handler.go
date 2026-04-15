package handler

import (
	"net/http"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gaston-garcia-cegid/arnela/backend/internal/service"
	"github.com/gaston-garcia-cegid/arnela/backend/pkg/cache"
	pkgerrors "github.com/gaston-garcia-cegid/arnela/backend/pkg/errors"
	"github.com/gin-gonic/gin"
)

// StatsHandler handles statistics-related endpoints
type StatsHandler struct {
	statsService *service.StatsService
	cache        *cache.CacheService
}

// NewStatsHandler creates a new StatsHandler
func NewStatsHandler(statsService *service.StatsService, cacheService ...*cache.CacheService) *StatsHandler {
	h := &StatsHandler{statsService: statsService}
	if len(cacheService) > 0 {
		h.cache = cacheService[0]
	}
	return h
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
	cacheKey := cache.DashboardStatsCacheKey()

	if h.cache != nil {
		var cached domain.DashboardStats
		if err := h.cache.Get(c.Request.Context(), cacheKey, &cached); err == nil {
			c.JSON(http.StatusOK, cached)
			return
		}
	}

	stats, err := h.statsService.GetDashboardStats(c.Request.Context())
	if err != nil {
		appErr := pkgerrors.NewInternalError("Error al obtener estadísticas")
		pkgerrors.RespondWithAppError(c, appErr)
		return
	}

	if h.cache != nil {
		_ = h.cache.Set(c.Request.Context(), cacheKey, stats, cache.CacheTTLShort)
	}

	c.JSON(http.StatusOK, stats)
}
