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

		// Cards routes (protected)
		cards := v1.Group("/cards")
		cards.Use(middleware.AuthMiddleware(jwtSecret))
		{
			cards.GET("", h.GetCards)
			cards.GET("/:id", h.GetCardByID)
			cards.POST("", h.CreateCard)
			cards.PUT("/:id", h.UpdateCard)
			cards.DELETE("/:id", h.DeleteCard)
		}

		// Card Owner routes (protected)
		cardOwners := v1.Group("/card-owners")
		cardOwners.Use(middleware.AuthMiddleware(jwtSecret))
		{
			cardOwners.POST("/register", h.RegisterCardOwner)          // Register single card
			cardOwners.POST("/register-multiple", h.RegisterMultipleCards) // Register multiple cards
			cardOwners.GET("/profile", h.GetCardOwnerProfile)          // Gets first card (backward compatibility)
			cardOwners.GET("/profiles", h.GetCardOwnerProfiles)        // Gets all cards for user
			cardOwners.PUT("/:id", h.UpdateCardOwner)                  // Update specific card owner by ID
			cardOwners.DELETE("/:id", h.DeleteCardOwner)               // Delete specific card owner by ID
			cardOwners.GET("/all", h.GetAllCardOwners)                 // Admin only
			
			// New API endpoints
			cardOwners.POST("/validate-duplicate", h.ValidateDuplicateCardRegistration) // Validate duplicate registration
			cardOwners.GET("/search/by-card", h.SearchCardOwnersByCardNameAndNumber)    // Search by card name and number
			cardOwners.GET("/search/by-owner", h.SearchCardOwnersByIDCardOrPhone)       // Search by ID card or phone
		}

		// Protected routes (example)
		protected := v1.Group("/protected")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			// Add protected endpoints here
		}
	}
}
