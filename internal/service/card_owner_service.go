package service

import (
	"errors"
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

// RegisterCardOwner creates a new card owner registration
func (s *CardOwnerService) RegisterCardOwner(userID uint, req *models.RegisterOwnerRequest) (*models.CardOwner, error) {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Validate that the card ID exists in master data
	_, err = s.repo.GetCardByID(req.CardID)
	if err != nil {
		return nil, errors.New("card ID not found in master data")
	}

	// Check if this card number for this card ID is already taken
	existingByCardNumberAndID, _ := s.repo.GetCardOwnerByCardNumberAndCardID(req.CardNumber, req.CardID)
	if existingByCardNumberAndID != nil {
		return nil, errors.New("card number is already registered for this card")
	}

	// Check if ID card is already taken
	existingByIDCard, _ := s.repo.GetCardOwnerByIDCard(req.IDCard)
	if existingByIDCard != nil {
		return nil, errors.New("ID card is already registered")
	}

	// Check if user already has a card owner registration
	existingByUser, _ := s.repo.GetCardOwnerByUserID(userID)
	if existingByUser != nil {
		return nil, errors.New("user already has a card owner registration")
	}

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

// GetCardOwnerProfile retrieves card owner profile with card details
func (s *CardOwnerService) GetCardOwnerProfile(userID uint) (*models.CardOwnerWithCard, error) {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Get card owner by user ID
	cardOwner, err := s.repo.GetCardOwnerByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get the associated card master data
	card, err := s.repo.GetCardByID(cardOwner.CardID)
	if err != nil {
		return nil, errors.New("associated card not found")
	}

	return &models.CardOwnerWithCard{
		CardOwner: *cardOwner,
		Card:      card,
	}, nil
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

// UpdateCardOwner updates card owner information
func (s *CardOwnerService) UpdateCardOwner(userID uint, req *models.UpdateCardOwnerRequest) (*models.CardOwner, error) {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Get existing card owner
	cardOwner, err := s.repo.GetCardOwnerByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.CardID != 0 {
		// Validate that the new card ID exists in master data
		_, err = s.repo.GetCardByID(req.CardID)
		if err != nil {
			return nil, errors.New("card ID not found in master data")
		}
		cardOwner.CardID = req.CardID
	}

	if req.CardNumber != "" {
		// Check if new card number for the card ID is already taken by another owner
		existingByCardNumberAndID, _ := s.repo.GetCardOwnerByCardNumberAndCardID(req.CardNumber, cardOwner.CardID)
		if existingByCardNumberAndID != nil && existingByCardNumberAndID.ID != cardOwner.ID {
			return nil, errors.New("card number is already registered for this card")
		}
		cardOwner.CardNumber = req.CardNumber
	}

	if req.IDCard != "" {
		// Check if new ID card is already taken by another owner
		existingByIDCard, _ := s.repo.GetCardOwnerByIDCard(req.IDCard)
		if existingByIDCard != nil && existingByIDCard.ID != cardOwner.ID {
			return nil, errors.New("ID card is already registered")
		}
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

// DeleteCardOwner deletes card owner registration
func (s *CardOwnerService) DeleteCardOwner(userID uint) error {
	// Check user authentication and active status
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return err
	}

	// Get existing card owner
	cardOwner, err := s.repo.GetCardOwnerByUserID(userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteCardOwner(cardOwner.ID)
}
