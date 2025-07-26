package handlers

import (
	"net/http"
	"tiger-fasttrack-card/internal/models"
	"tiger-fasttrack-card/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{
		Service: svc,
	}
}

// Health check handler
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Tiger FastTrack Card API is running",
	})
}

// Example handlers for different endpoints

// GetCards handler
func (h *Handler) GetCards(c *gin.Context) {
	// TODO: Implement get cards logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Get cards endpoint",
		"data":    []interface{}{},
	})
}

// GetCardByID handler
func (h *Handler) GetCardByID(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement get card by ID logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Get card by ID endpoint",
		"id":      id,
	})
}

// CreateCard handler
func (h *Handler) CreateCard(c *gin.Context) {
	// TODO: Implement create card logic
	c.JSON(http.StatusCreated, gin.H{
		"message": "Create card endpoint",
	})
}

// UpdateCard handler
func (h *Handler) UpdateCard(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement update card logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Update card endpoint",
		"id":      id,
	})
}

// DeleteCard handler
func (h *Handler) DeleteCard(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement delete card logic
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete card endpoint",
		"id":      id,
	})
}

// Authentication handlers

// Register handler
func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}

// Login handler
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.Service.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handler
func (h *Handler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newToken, err := h.Service.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": newToken,
	})
}

// GetProfile handler
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.Service.GetUserProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdateProfile handler
func (h *Handler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Service.UpdateUserProfile(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

// ChangePassword handler
func (h *Handler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.ChangePassword(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully",
	})
}
