package service

import (
	"errors"
	"strings"
	"tiger-fasttrack-card/internal/models"
	"tiger-fasttrack-card/internal/repository"
)

// CardOwnerService handles card owner operations
type CardOwnerService struct {
	repo        *repository.Repository
	authService *AuthService
}

// NewCardOwnerService creates a new CardOwnerService instance
func NewCardOwnerService(repo *repository.Repository, authService *AuthService) *CardOwnerService {
	return &CardOwnerService{
		repo:        repo,
		authService: authService,
	}
}

// validateDuplicateCardRegistration checks if a card registration already exists
// excludeID can be provided to exclude a specific card owner ID from the duplicate check (useful for updates)
func (s *CardOwnerService) validateDuplicateCardRegistration(cardID uint, cardNumber string, excludeID ...uint) error {
	// Validate that the card ID exists in master data
	_, err := s.repo.GetCardByID(cardID)
	if err != nil {
		return errors.New("card ID not found in master data")
	}

	// Check if this card number for this card ID is already taken
	existingByCardNumberAndID, _ := s.repo.GetCardOwnerByCardNumberAndCardID(cardNumber, cardID)
	if existingByCardNumberAndID != nil {
		// If excludeID is provided, check if the existing record is the one being excluded
		if len(excludeID) > 0 && existingByCardNumberAndID.ID == excludeID[0] {
			return nil // This is an update of the same record, allow it
		}
		return errors.New("card number " + cardNumber + " is already registered for this card")
	}

	return nil
}

// RegisterCardOwner creates a new card owner registration
func (s *CardOwnerService) RegisterCardOwner(userID uint, req *models.RegisterOwnerRequest) (*models.CardOwner, error) {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Validate duplicate card registration
	err = s.validateDuplicateCardRegistration(req.CardID, req.CardNumber)
	if err != nil {
		return nil, err
	}

	// Note: Same owner (ID card) can register multiple different cards
	// We only prevent duplicate card_number + card_id combinations, not duplicate ID cards

	// Create card owner
	cardOwner := &models.CardOwner{
		CardID:      req.CardID,
		CardNumber:  req.CardNumber,
		IDCard:      req.IDCard,
		PhoneNumber: req.PhoneNumber,
		UserID:      userID,
	}

	err = s.repo.CreateCardOwner(cardOwner)
	if err != nil {
		return nil, errors.New("failed to register card owner")
	}

	return cardOwner, nil
}

// RegisterMultipleCards creates multiple card owner registrations in a single transaction
func (s *CardOwnerService) RegisterMultipleCards(userID uint, req *models.RegisterMultipleCardsRequest) ([]models.CardOwner, error) {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	var cardOwners []models.CardOwner
	
	// Validate each card and check for duplicates
	for _, cardReg := range req.Cards {
		err = s.validateDuplicateCardRegistration(cardReg.CardID, cardReg.CardNumber)
		if err != nil {
			return nil, err
		}
	}

	// Create all card owner registrations
	for _, cardReg := range req.Cards {
		cardOwner := &models.CardOwner{
			CardID:      cardReg.CardID,
			CardNumber:  cardReg.CardNumber,
			IDCard:      req.IDCard,
			PhoneNumber: req.PhoneNumber,
			UserID:      userID,
		}

		err = s.repo.CreateCardOwner(cardOwner)
		if err != nil {
			return nil, errors.New("failed to register card owner")
		}

		cardOwners = append(cardOwners, *cardOwner)
	}

	return cardOwners, nil
}

// GetCardOwnerProfiles retrieves all card owner profiles for a user with card details
func (s *CardOwnerService) GetCardOwnerProfiles(userID uint) ([]models.CardOwnerWithCard, error) {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Get all card owners by user ID
	cardOwners, err := s.repo.GetCardOwnersByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Populate card master data for each owner
	var result []models.CardOwnerWithCard
	for _, owner := range cardOwners {
		card, err := s.repo.GetCardByID(owner.CardID)
		if err != nil {
			// Skip if card not found, but continue processing others
			continue
		}
		result = append(result, models.CardOwnerWithCard{
			CardOwner: owner,
			Card:      card,
		})
	}

	return result, nil
}

// GetAllCardOwners retrieves all card owners (admin only)
func (s *CardOwnerService) GetAllCardOwners(userID uint) ([]models.CardOwnerWithCard, error) {
	// Check user authentication and active status (admin only)
	user, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}
	if user.Role != "admin" {
		return nil, errors.New("insufficient permissions")
	}

	// Get all card owners
	cardOwners, err := s.repo.GetAllCardOwners()
	if err != nil {
		return nil, errors.New("failed to get card owners")
	}

	// Populate card master data for each owner
	var result []models.CardOwnerWithCard
	for _, owner := range cardOwners {
		card, err := s.repo.GetCardByID(owner.CardID)
		if err != nil {
			// Skip if card not found, but continue processing others
			continue
		}
		result = append(result, models.CardOwnerWithCard{
			CardOwner: owner,
			Card:      card,
		})
	}

	return result, nil
}

// UpdateCardOwner updates card owner information by card owner ID
func (s *CardOwnerService) UpdateCardOwner(userID uint, cardOwnerID uint, req *models.UpdateCardOwnerRequest) (*models.CardOwner, error) {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Get existing card owner by ID and verify ownership
	cardOwner, err := s.repo.GetCardOwnerByID(cardOwnerID)
	if err != nil {
		return nil, err
	}

	// Verify that the card owner belongs to the authenticated user
	if cardOwner.UserID != userID {
		return nil, errors.New("card owner not found or access denied")
	}

	// Update fields if provided
	if req.CardID != 0 {
		cardOwner.CardID = req.CardID
	}

	if req.CardNumber != "" {
		cardOwner.CardNumber = req.CardNumber
	}

	// If either CardID or CardNumber is being updated, validate the combination
	if req.CardID != 0 || req.CardNumber != "" {
		err = s.validateDuplicateCardRegistration(cardOwner.CardID, cardOwner.CardNumber, cardOwner.ID)
		if err != nil {
			return nil, err
		}
	}

	if req.IDCard != "" {
		// Allow updating ID card (same owner can have multiple registrations)
		cardOwner.IDCard = req.IDCard
	}

	if req.PhoneNumber != "" {
		cardOwner.PhoneNumber = req.PhoneNumber
	}

	err = s.repo.UpdateCardOwner(cardOwner)
	if err != nil {
		return nil, errors.New("failed to update card owner")
	}

	return cardOwner, nil
}

// DeleteCardOwner deletes a specific card owner registration by card owner ID
func (s *CardOwnerService) DeleteCardOwner(userID uint, cardOwnerID uint) error {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return err
	}

	// Get existing card owner by ID and verify ownership
	cardOwner, err := s.repo.GetCardOwnerByID(cardOwnerID)
	if err != nil {
		return err
	}

	// Verify that the card owner belongs to the authenticated user
	if cardOwner.UserID != userID {
		return errors.New("card owner not found or access denied")
	}

	return s.repo.DeleteCardOwner(cardOwner.ID)
}

// ValidateDuplicateCardRegistration is a public service for validating duplicate card registration
// This can be used by external services or handlers to check for duplicates before registration
func (s *CardOwnerService) ValidateDuplicateCardRegistration(cardID uint, cardNumber string) error {
	return s.validateDuplicateCardRegistration(cardID, cardNumber)
}

// ValidateDuplicateCardRegistrationForUpdate is a public service for validating duplicate card registration during updates
// excludeID should be the ID of the card owner being updated
func (s *CardOwnerService) ValidateDuplicateCardRegistrationForUpdate(cardID uint, cardNumber string, excludeID uint) error {
	return s.validateDuplicateCardRegistration(cardID, cardNumber, excludeID)
}

// SearchCardOwnersByCardNameAndNumber searches for card owners by card name and card number
// This service allows searching across card registrations using card master data
func (s *CardOwnerService) SearchCardOwnersByCardNameAndNumber(userID uint, cardName string, cardNumber string) ([]models.CardOwnerWithCard, error) {
	// Check user authentication and active status
	user, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	var result []models.CardOwnerWithCard
	
	// If admin, search across all card owners; if regular user, search only their own registrations
	var cardOwners []models.CardOwner
	if user.Role == "admin" {
		cardOwners, err = s.repo.GetAllCardOwners()
		if err != nil {
			return nil, errors.New("failed to get card owners")
		}
	} else {
		cardOwners, err = s.repo.GetCardOwnersByUserID(userID)
		if err != nil {
			return nil, err
		}
	}

	// Filter by card name and card number
	for _, owner := range cardOwners {
		// Get card details
		card, err := s.repo.GetCardByID(owner.CardID)
		if err != nil {
			continue // Skip if card not found
		}

		// Check if card name matches (case-insensitive partial match)
		cardNameMatches := cardName == "" || 
			len(cardName) == 0 || 
			strings.Contains(strings.ToLower(card.CardName), strings.ToLower(cardName))

		// Check if card number matches (exact match or partial match)
		cardNumberMatches := cardNumber == "" || 
			len(cardNumber) == 0 || 
			strings.Contains(owner.CardNumber, cardNumber)

		if cardNameMatches && cardNumberMatches {
			result = append(result, models.CardOwnerWithCard{
				CardOwner: owner,
				Card:      card,
			})
		}
	}

	return result, nil
}

// SearchCardOwnersByIDCardOrPhone searches for card owners by ID card or phone number
// This service allows searching for card owners using their personal information
func (s *CardOwnerService) SearchCardOwnersByIDCardOrPhone(userID uint, idCard string, phoneNumber string) ([]models.CardOwnerWithCard, error) {
	// Check user authentication and active status
	user, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// At least one search parameter must be provided
	if (idCard == "" || len(idCard) == 0) && (phoneNumber == "" || len(phoneNumber) == 0) {
		return nil, errors.New("at least one search parameter (ID card or phone number) must be provided")
	}

	var result []models.CardOwnerWithCard
	
	// If admin, search across all card owners; if regular user, search only their own registrations
	var cardOwners []models.CardOwner
	if user.Role == "admin" {
		cardOwners, err = s.repo.GetAllCardOwners()
		if err != nil {
			return nil, errors.New("failed to get card owners")
		}
	} else {
		cardOwners, err = s.repo.GetCardOwnersByUserID(userID)
		if err != nil {
			return nil, err
		}
	}

	// Filter by ID card or phone number
	for _, owner := range cardOwners {
		// Check if ID card matches (partial match, case-insensitive)
		idCardMatches := idCard == "" || 
			len(idCard) == 0 || 
			strings.Contains(strings.ToLower(owner.IDCard), strings.ToLower(idCard))

		// Check if phone number matches (partial match)
		phoneMatches := phoneNumber == "" || 
			len(phoneNumber) == 0 || 
			strings.Contains(owner.PhoneNumber, phoneNumber)

		// Match if either ID card OR phone number matches (OR logic)
		if idCardMatches || phoneMatches {
			// Get card details
			card, err := s.repo.GetCardByID(owner.CardID)
			if err != nil {
				continue // Skip if card not found
			}

			result = append(result, models.CardOwnerWithCard{
				CardOwner: owner,
				Card:      card,
			})
		}
	}

	return result, nil
}


