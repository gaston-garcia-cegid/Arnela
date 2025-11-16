package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

// RedisClient wraps the Redis client
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient creates a new Redis connection
func NewRedisClient(addr, password string, db int) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error connecting to Redis: %w", err)
	}

	log.Println("✓ Redis connected successfully")

	return &RedisClient{Client: client}, nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.Client.Close()
}
