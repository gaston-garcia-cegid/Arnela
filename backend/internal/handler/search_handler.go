package handler

import (
	"net/http"
	"strings"

	"github.com/gaston-garcia-cegid/arnela/backend/internal/domain"
	"github.com/gin-gonic/gin"
)

// SearchHandler handles global search requests
type SearchHandler struct {
	searchService domain.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService domain.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// GlobalSearch godoc
// @Summary Global search across all entities
// @Description Search for clients, employees, appointments, and invoices by query string
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query (minimum 2 characters)"
// @Success 200 {object} domain.SearchResults
// @Failure 400 {object} ErrorResponse "Invalid query parameter"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /search [get]
func (h *SearchHandler) GlobalSearch(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))

	// Validate query parameter
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query parameter 'q' is required",
		})
		return
	}

	// Minimum 2 characters for search
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "query must be at least 2 characters long",
		})
		return
	}

	// Perform search
	results, err := h.searchService.GlobalSearch(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to perform search: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, results)
}
