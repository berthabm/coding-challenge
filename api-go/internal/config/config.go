// Package config centralizes environment-driven settings.
// Keeps handlers/services free of os.Getenv calls (testability + single source of truth).
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for API 1.
type Config struct {
	AppName  string
	Env      string
	Host     string
	Port     string
	LogLevel string

	API2BaseURL    string
	API2Timeout    time.Duration
	API2MatrixPath string

	// JWT auth
	AuthUsername string
	AuthPassword string
	JWTSecret    string

	// CORS
	CORSAllowedOrigins string
}

// Load reads .env (optional in production) and environment variables with defaults.
func Load() (*Config, error) {
	_ = godotenv.Load()

	timeoutSec, err := strconv.Atoi(getEnv("API2_TIMEOUT_SECONDS", "10"))
	if err != nil {
		return nil, fmt.Errorf("API2_TIMEOUT_SECONDS: %w", err)
	}

	return &Config{
		AppName:        getEnv("APP_NAME", "api-go"),
		Env:            getEnv("APP_ENV", "development"),
		Host:           getEnv("APP_HOST", "0.0.0.0"),
		Port:           getEnv("APP_PORT", "8080"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		API2BaseURL:    getEnv("API2_BASE_URL", "http://localhost:3001"),
		API2Timeout:    time.Duration(timeoutSec) * time.Second,
		API2MatrixPath: getEnv("API2_MATRIX_PATH", "/api/stats"),
		AuthUsername:        getEnv("AUTH_USERNAME", "admin"),
		AuthPassword:        getEnv("AUTH_PASSWORD", "admin123"),
		JWTSecret:           getEnv("JWT_SECRET", "change-me-in-production"),
		CORSAllowedOrigins:  getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
