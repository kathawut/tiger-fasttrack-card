package handlers

import (
	"net/http"
	"strconv"
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
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	cards, err := h.Service.GetAllCards(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cards retrieved successfully",
		"data":    cards,
	})
}

// GetCardByID handler
func (h *Handler) GetCardByID(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	card, err := h.Service.GetCardByID(userID.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card retrieved successfully",
		"data":    card,
	})
}

// CreateCard handler
func (h *Handler) CreateCard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.Service.CreateCard(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Card created successfully",
		"data":    card,
	})
}

// UpdateCard handler
func (h *Handler) UpdateCard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	var req models.UpdateCardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	card, err := h.Service.UpdateCard(userID.(uint), uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card updated successfully",
		"data":    card,
	})
}

// DeleteCard handler
func (h *Handler) DeleteCard(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card ID"})
		return
	}

	err = h.Service.DeleteCard(userID.(uint), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card deleted successfully",
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

// CardOwner handlers

// RegisterCardOwner handler
func (h *Handler) RegisterCardOwner(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.RegisterOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardOwner, err := h.Service.RegisterCardOwner(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Card owner registered successfully",
		"data":    cardOwner,
	})
}

// RegisterMultipleCards handler
func (h *Handler) RegisterMultipleCards(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.RegisterMultipleCardsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardOwners, err := h.Service.RegisterMultipleCards(userID.(uint), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Multiple cards registered successfully",
		"data":    cardOwners,
	})
}

// GetCardOwnerProfile handler
func (h *Handler) GetCardOwnerProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profile, err := h.Service.GetCardOwnerProfile(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card owner profile retrieved successfully",
		"data":    profile,
	})
}

// GetCardOwnerProfiles handler - gets all card registrations for authenticated user
func (h *Handler) GetCardOwnerProfiles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	profiles, err := h.Service.GetCardOwnerProfiles(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card owner profiles retrieved successfully",
		"data":    profiles,
	})
}

// GetAllCardOwners handler (admin only)
func (h *Handler) GetAllCardOwners(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	cardOwners, err := h.Service.GetAllCardOwners(userID.(uint))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card owners retrieved successfully",
		"data":    cardOwners,
	})
}

// UpdateCardOwner handler
func (h *Handler) UpdateCardOwner(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get card owner ID from URL parameter
	cardOwnerIDStr := c.Param("id")
	cardOwnerID, err := strconv.ParseUint(cardOwnerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card owner ID"})
		return
	}

	var req models.UpdateCardOwnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cardOwner, err := h.Service.UpdateCardOwner(userID.(uint), uint(cardOwnerID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card owner updated successfully",
		"data":    cardOwner,
	})
}

// DeleteCardOwner handler
func (h *Handler) DeleteCardOwner(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get card owner ID from URL parameter
	cardOwnerIDStr := c.Param("id")
	cardOwnerID, err := strconv.ParseUint(cardOwnerIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid card owner ID"})
		return
	}

	err = h.Service.DeleteCardOwner(userID.(uint), uint(cardOwnerID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Card owner deleted successfully",
	})
}

// ValidateDuplicateCardRegistration validates if a card registration would be duplicate
func (h *Handler) ValidateDuplicateCardRegistration(c *gin.Context) {
	var req struct {
		CardID     uint   `json:"card_id" binding:"required"`
		CardNumber string `json:"card_number" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.Service.CardOwnerService.ValidateDuplicateCardRegistration(req.CardID, req.CardNumber)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":     "Duplicate card registration",
			"message":   err.Error(),
			"duplicate": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Card registration is valid",
		"duplicate": false,
	})
}

// SearchCardOwnersByCardNameAndNumber searches for card owners by card name and card number
func (h *Handler) SearchCardOwnersByCardNameAndNumber(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get query parameters
	cardName := c.Query("card_name")
	cardNumber := c.Query("card_number")

	// At least one parameter should be provided
	if cardName == "" && cardNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one search parameter (card_name or card_number) must be provided",
		})
		return
	}

	cardOwners, err := h.Service.CardOwnerService.SearchCardOwnersByCardNameAndNumber(userID.(uint), cardName, cardNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Search completed successfully",
		"data":    cardOwners,
		"count":   len(cardOwners),
	})
}

// SearchCardOwnersByIDCardOrPhone searches for card owners by ID card or phone number
func (h *Handler) SearchCardOwnersByIDCardOrPhone(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get query parameters
	idCard := c.Query("id_card")
	phoneNumber := c.Query("phone_number")

	cardOwners, err := h.Service.CardOwnerService.SearchCardOwnersByIDCardOrPhone(userID.(uint), idCard, phoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Search completed successfully",
		"data":    cardOwners,
		"count":   len(cardOwners),
	})
}
