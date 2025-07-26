package routes

import (
	"tiger-fasttrack-card/internal/handlers"
	"tiger-fasttrack-card/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine, h *handlers.Handler, jwtSecret string) {
	// Health check endpoint
	router.GET("/health", h.HealthCheck)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", h.Register)
			auth.POST("/login", h.Login)
			auth.POST("/refresh", h.RefreshToken)
		}

		// User routes
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(jwtSecret))
		{
			users.GET("/profile", h.GetProfile)
			users.PUT("/profile", h.UpdateProfile)
			users.POST("/change-password", h.ChangePassword)
		}

		// Cards routes
		cards := v1.Group("/cards")
		{
			cards.GET("", h.GetCards)
			cards.GET("/:id", h.GetCardByID)
			cards.POST("", h.CreateCard)
			cards.PUT("/:id", h.UpdateCard)
			cards.DELETE("/:id", h.DeleteCard)
		}

		// Protected routes (example)
		protected := v1.Group("/protected")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Add protected endpoints here
		}
	}
}
