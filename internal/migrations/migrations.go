package migrations

import (
	"tiger-fasttrack-card/internal/database"
	"tiger-fasttrack-card/internal/models"
)

// RunMigrations executes all database migrations
func RunMigrations(db *database.Database) error {
	// Add your models here when you create them
	return db.GetDB().AutoMigrate(
		&models.User{},
		// Add other models here as you create them
		// &models.Card{},
		// &models.Transaction{},
	)
}
