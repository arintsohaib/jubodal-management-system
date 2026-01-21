package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	// Database
	DatabaseURL string

	// Redis
	RedisURL      string
	RedisPassword string

	// OpenSearch
	OpenSearchURL string

	// JWT
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration

	// Server
	Port           string
	Environment    string
	AllowedOrigins []string

	// Rate Limiting
	RateLimitLogin  int
	RateLimitWindow time.Duration

	// Security
	MaxFailedAttempts    int
	LockoutDuration      time.Duration
	MaxConcurrentSessions int

	// Logging
	LogLevel  string
	LogFormat string
}

// Load loads configuration from environment variables
func Load() *Config {
	// Load .env file if exists (development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		RedisURL:             getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		OpenSearchURL:        getEnv("OPENSEARCH_URL", "http://opensearch:9200"),
		JWTAccessSecret:      getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret:     getEnv("JWT_REFRESH_SECRET", ""),
		JWTAccessExpiry:      getDuration("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry:     getDuration("JWT_REFRESH_EXPIRY", "2160h"), // 90 days
		Port:                 getEnv("PORT", "3000"),
		Environment:          getEnv("ENV", "development"),
		RateLimitLogin:       10,
		RateLimitWindow:      15 * time.Minute,
		MaxFailedAttempts:    5,
		LockoutDuration:      30 * time.Minute,
		MaxConcurrentSessions: 3,
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		LogFormat:            getEnv("LOG_FORMAT", "json"),
	}
}

// Validate checks if required configuration is present
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	if c.JWTAccessSecret == "" {
		log.Fatal("JWT_ACCESS_SECRET is required")
	}
	if c.JWTRefreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET is required")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDuration(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("Invalid duration for %s: %v", key, err)
	}
	return duration
}
