package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheService provides caching functionality
type CacheService struct {
	client *redis.Client
}

// NewCacheService creates a new cache service
func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{
		client: client,
	}
}

// Get retrieves a value from cache and unmarshals it into dest
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Set stores a value in cache with expiration
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return c.client.Set(ctx, key, data, expiration).Err()
}

// Delete removes a value from cache
func (c *CacheService) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// DeletePattern deletes all keys matching a pattern
func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, pattern, 0).Iterator()

	keys := []string{}
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return c.client.Del(ctx, keys...).Err()
	}

	return nil
}

// Exists checks if a key exists in cache
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetOrSet retrieves from cache or sets if not exists (Cache-Aside pattern)
func (c *CacheService) GetOrSet(ctx context.Context, key string, dest interface{}, expiration time.Duration, loader func() (interface{}, error)) error {
	// Try to get from cache
	err := c.Get(ctx, key, dest)
	if err == nil {
		// Cache hit
		return nil
	}

	if err != redis.Nil {
		// Real error (not just cache miss)
		return err
	}

	// Cache miss - load from source
	data, err := loader()
	if err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	// Store in cache
	if err := c.Set(ctx, key, data, expiration); err != nil {
		// Log error but don't fail the request
		// The data is still valid even if caching fails
		return nil
	}

	// Marshal into dest
	jsonData, _ := json.Marshal(data)
	return json.Unmarshal(jsonData, dest)
}

// Cache key generators for different entities

// ClientCacheKey generates a cache key for a client
func ClientCacheKey(clientID int64) string {
	return fmt.Sprintf("client:%d", clientID)
}

// ClientListCacheKey generates a cache key for client list
func ClientListCacheKey() string {
	return "clients:list"
}

// EmployeeCacheKey generates a cache key for an employee
func EmployeeCacheKey(employeeID int64) string {
	return fmt.Sprintf("employee:%d", employeeID)
}

// EmployeeListCacheKey generates a cache key for employee list
func EmployeeListCacheKey() string {
	return "employees:list"
}

// EmployeesBySpecialtyCacheKey generates a cache key for employees by specialty
func EmployeesBySpecialtyCacheKey(specialty string) string {
	return fmt.Sprintf("employees:specialty:%s", specialty)
}

// AppointmentCacheKey generates a cache key for an appointment
func AppointmentCacheKey(appointmentID int64) string {
	return fmt.Sprintf("appointment:%d", appointmentID)
}

// ClientAppointmentsCacheKey generates a cache key for client appointments
func ClientAppointmentsCacheKey(clientID int64) string {
	return fmt.Sprintf("appointments:client:%d", clientID)
}

// EmployeeAppointmentsCacheKey generates a cache key for employee appointments
func EmployeeAppointmentsCacheKey(employeeID int64, date string) string {
	return fmt.Sprintf("appointments:employee:%d:date:%s", employeeID, date)
}

// DashboardStatsCacheKey generates a cache key for dashboard stats
func DashboardStatsCacheKey() string {
	return "stats:dashboard"
}

// Cache TTL constants
const (
	CacheTTLShort  = 5 * time.Minute  // For frequently changing data
	CacheTTLMedium = 15 * time.Minute // For moderately changing data
	CacheTTLLong   = 1 * time.Hour    // For rarely changing data
	CacheTTLDay    = 24 * time.Hour   // For very stable data
)
