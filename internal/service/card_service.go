package service

import (
	"errors"
	"tiger-fasttrack-card/internal/models"
	"tiger-fasttrack-card/internal/repository"
)

type CardService struct {
	repo        *repository.Repository
	authService *AuthService
}

func NewCardService(repo *repository.Repository, authService *AuthService) *CardService {
	return &CardService{
		repo:        repo,
		authService: authService,
	}
}

// Card service methods
func (s *CardService) GetAllCards(userID uint) ([]models.Card, error) {
	// Verify user exists and is active
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAllCards()
}

func (s *CardService) GetCardByID(userID uint, id uint) (*models.Card, error) {
	// Verify user exists and is active
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetCardByID(id)
}

func (s *CardService) CreateCard(userID uint, req *models.CreateCardRequest) (*models.Card, error) {
	// Verify user exists and is active
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Validate card quantity
	if req.CardQuantity < 0 {
		return nil, errors.New("card quantity cannot be negative")
	}

	// Create card
	card := &models.Card{
		CardName:     req.CardName,
		CardImage:    req.CardImage,
		CardQuantity: req.CardQuantity,
	}

	err = s.repo.CreateCard(card)
	if err != nil {
		return nil, errors.New("failed to create card")
	}

	return card, nil
}

func (s *CardService) UpdateCard(userID uint, id uint, req *models.UpdateCardRequest) (*models.Card, error) {
	// Verify user exists and is active
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return nil, err
	}

	// Get existing card
	card, err := s.repo.GetCardByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.CardName != "" {
		card.CardName = req.CardName
	}
	if req.CardImage != "" {
		card.CardImage = req.CardImage
	}
	if req.CardQuantity != nil {
		if *req.CardQuantity < 0 {
			return nil, errors.New("card quantity cannot be negative")
		}
		card.CardQuantity = *req.CardQuantity
	}

	err = s.repo.UpdateCard(card)
	if err != nil {
		return nil, errors.New("failed to update card")
	}

	return card, nil
}

func (s *CardService) DeleteCard(userID uint, id uint) error {
	// Verify user exists and is active
	_, err := s.authService.ValidateUserAccess(userID)
	if err != nil {
		return err
	}

	// Check if card exists
	_, err = s.repo.GetCardByID(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteCard(id)
}
