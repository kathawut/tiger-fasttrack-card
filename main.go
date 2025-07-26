package main

import (
	"log"
	"os"

	"tiger-fasttrack-card/internal/config"
	"tiger-fasttrack-card/internal/database"
	"tiger-fasttrack-card/internal/handlers"
	"tiger-fasttrack-card/internal/middleware"
	"tiger-fasttrack-card/internal/migrations"
	"tiger-fasttrack-card/internal/repository"
	"tiger-fasttrack-card/internal/routes"
	"tiger-fasttrack-card/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize configuration
	cfg := config.New()

	// Initialize database
	db, err := database.New(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run database migrations
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// Initialize repository
	repo := repository.New(db)

	// Initialize service with JWT secret
	svc := service.New(repo, cfg.JWTSecret)

	// Initialize handlers
	h := handlers.New(svc)

	// Setup routes
	routes.Setup(router, h, cfg.JWTSecret)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
