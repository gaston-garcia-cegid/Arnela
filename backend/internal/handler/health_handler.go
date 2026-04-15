package handler

import (
	"context"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

var startTime = time.Now()

type HealthHandler struct {
	db          *sqlx.DB
	redisClient *redis.Client
	version     string
}

func NewHealthHandler(db *sqlx.DB, redisClient *redis.Client, version string) *HealthHandler {
	return &HealthHandler{db: db, redisClient: redisClient, version: version}
}

func (h *HealthHandler) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
	defer cancel()

	dbStatus := "connected"
	if err := h.db.PingContext(ctx); err != nil {
		dbStatus = "disconnected"
	}

	redisStatus := "connected"
	if err := h.redisClient.Ping(ctx).Err(); err != nil {
		redisStatus = "disconnected"
	}

	healthy := dbStatus == "connected" && redisStatus == "connected"
	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	dbStats := h.db.Stats()

	c.JSON(status, gin.H{
		"status":  map[bool]string{true: "healthy", false: "unhealthy"}[healthy],
		"version": h.version,
		"uptime":  time.Since(startTime).String(),
		"services": gin.H{
			"database": gin.H{
				"status":      dbStatus,
				"connections": dbStats.OpenConnections,
				"idle":        dbStats.Idle,
				"in_use":      dbStats.InUse,
			},
			"redis": gin.H{
				"status": redisStatus,
			},
		},
		"runtime": gin.H{
			"go_version":   runtime.Version(),
			"goroutines":   runtime.NumGoroutine(),
			"alloc_mb":     mem.Alloc / 1024 / 1024,
			"sys_mb":       mem.Sys / 1024 / 1024,
			"gc_completed": mem.NumGC,
		},
	})
}

// Readiness is a lightweight probe for orchestrators (K8s, Docker healthcheck).
func (h *HealthHandler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"ready": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"ready": true})
}
