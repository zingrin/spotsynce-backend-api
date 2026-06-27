package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	Port           string
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	JWTSecret      string
	JWTExpiryHours int
}

// Load reads configuration from .env file (if present) and environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiry := 24
	if v := os.Getenv("JWT_EXPIRY_HOURS"); v != "" {
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
		}
		jwtExpiry = parsed
	}

	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBName:         getEnv("DB_NAME", "spotsync"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiryHours: jwtExpiry,
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=UTC",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

// JWTExpiryDuration returns the JWT token lifetime.
func (c *Config) JWTExpiryDuration() time.Duration {
	return time.Duration(c.JWTExpiryHours) * time.Hour
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
