package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// AppConfig holds application-level configuration
type AppConfig struct {
	Env  string
	Port string
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "arnela_user"),
			Password: getEnv("DB_PASSWORD", "arnela_secure_pass_2024"),
			DBName:   getEnv("DB_NAME", "arnela_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", "arnela_redis_pass_2024"),
			DB:       0,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your_jwt_secret_key_change_in_production"),
		},
	}
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// GetRedisAddress returns the Redis connection address
func (c *RedisConfig) GetRedisAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// getEnv gets an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
