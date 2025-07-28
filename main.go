package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tiger-fasttrack-card/internal/config"
	"tiger-fasttrack-card/internal/database"
	"tiger-fasttrack-card/internal/handlers"
	"tiger-fasttrack-card/internal/middleware"
	"tiger-fasttrack-card/internal/migrations"
	"tiger-fasttrack-card/internal/repository"
	"tiger-fasttrack-card/internal/routes"
	"tiger-fasttrack-card/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Note: Using only environment variables from DigitalOcean App Platform
	// No .env file loading needed in production

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

	// Start server with graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = cfg.Port
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
