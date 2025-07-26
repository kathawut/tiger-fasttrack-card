package database

import (
	"fmt"
	"strings"
	"tiger-fasttrack-card/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func New(cfg *config.Config) (*Database, error) {
	var dsn string
	
	// Debug: print configuration values
	fmt.Printf("DEBUG: DatabaseURL: '%s'\n", cfg.DatabaseURL)
	fmt.Printf("DEBUG: DB_HOST: '%s'\n", cfg.Database.Host)
	fmt.Printf("DEBUG: DB_USER: '%s'\n", cfg.Database.User)
	fmt.Printf("DEBUG: DB_NAME: '%s'\n", cfg.Database.Name)
	
	// Use DATABASE_URL only if it's a proper PostgreSQL URL
	if cfg.DatabaseURL != "" && (strings.HasPrefix(cfg.DatabaseURL, "postgres://") || strings.HasPrefix(cfg.DatabaseURL, "postgresql://")) {
		dsn = cfg.DatabaseURL
		fmt.Printf("DEBUG: Using DATABASE_URL\n")
	} else {
		// Build DSN from individual components
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			cfg.Database.Host,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name,
			cfg.Database.Port,
			cfg.Database.SSLMode,
		)
		fmt.Printf("DEBUG: Using individual components. DSN: %s\n", dsn)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) Migrate() error {
	// Import and use migrations package
	return d.DB.AutoMigrate(
		// Add models here or import from migrations package
	)
}

func (d *Database) GetDB() *gorm.DB {
	return d.DB
}
