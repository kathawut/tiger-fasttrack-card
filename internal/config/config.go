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
		DatabaseURL: getEnv("DATABASE_URL", "postgresql://tiger-app:AVNS_c5wTtwqDbeQ2wWTmcTo@private-tiger-card-db-do-user-527620-0.g.db.ondigitalocean.com:25060/tiger-card?sslmode=verify-full"),
		JWTSecret:   getEnv("JWT_SECRET", "79dbfcfcd07688b791aacd22fb797354"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "private-tiger-card-db-do-user-527620-0.g.db.ondigitalocean.com"),
			Port:     getEnv("DB_PORT", "25060"),
			User:     getEnv("DB_USER", "tiger-app"),
			Password: getEnv("DB_PASSWORD", "AVNS_c5wTtwqDbeQ2wWTmcTo"),
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
