package config

import (
	"os"
)

type Config struct {
	Environment string
	Port        string
	DatabaseURL string
	JWTSecret   string
	Database    DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func New() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "production"),
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://doadmin:your-password@private-tiger-card-db-do-user-527620-0.g.db.ondigitalocean.com:25060/tiger-card?sslmode=require"),
		JWTSecret:   getEnv("JWT_SECRET", "your-production-jwt-secret"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "private-tiger-card-db-do-user-527620-0.g.db.ondigitalocean.com"),
			Port:     getEnv("DB_PORT", "25060"),
			User:     getEnv("DB_USER", "doadmin"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "tiger-card"),
			SSLMode:  getEnv("DB_SSLMODE", "verify-full"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
