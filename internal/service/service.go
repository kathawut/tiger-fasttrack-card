package service

import (
	"errors"
	"tiger-fasttrack-card/internal/repository"
	"tiger-fasttrack-card/internal/utils"
	"tiger-fasttrack-card/internal/models"
)

// Service is the main service that coordinates all sub-services
type Service struct {
	AuthService      *AuthService
	CardService      *CardService
	CardOwnerService *CardOwnerService
}

// New creates a new Service instance with all sub-services
func New(repo *repository.Repository, jwtSecret string) *Service {
	jwtManager := utils.NewJWTManager(jwtSecret)
	
	// Create auth service first as other services depend on it
	authService := NewAuthService(repo, jwtManager)
	
	// Create other services with auth service dependency
	cardService := NewCardService(repo, authService)
	cardOwnerService := NewCardOwnerService(repo, authService)

	return &Service{
		AuthService:      authService,
		CardService:      cardService,
		CardOwnerService: cardOwnerService,
	}
}

// Authentication service delegation methods
func (s *Service) Register(req *models.RegisterRequest) (*models.User, error) {
	return s.AuthService.Register(req)
}

func (s *Service) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	return s.AuthService.Login(req)
}

func (s *Service) RefreshToken(refreshToken string) (string, error) {
	return s.AuthService.RefreshToken(refreshToken)
}

func (s *Service) GetUserProfile(userID uint) (*models.User, error) {
	return s.AuthService.GetUserProfile(userID)
}

func (s *Service) UpdateUserProfile(userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	return s.AuthService.UpdateUserProfile(userID, req)
}

func (s *Service) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	return s.AuthService.ChangePassword(userID, req)
}

// Card service delegation methods
func (s *Service) GetAllCards(userID uint) ([]models.Card, error) {
	return s.CardService.GetAllCards(userID)
}

func (s *Service) GetCardByID(userID uint, id uint) (*models.Card, error) {
	return s.CardService.GetCardByID(userID, id)
}

func (s *Service) CreateCard(userID uint, req *models.CreateCardRequest) (*models.Card, error) {
	return s.CardService.CreateCard(userID, req)
}

func (s *Service) UpdateCard(userID uint, id uint, req *models.UpdateCardRequest) (*models.Card, error) {
	return s.CardService.UpdateCard(userID, id, req)
}

func (s *Service) DeleteCard(userID uint, id uint) error {
	return s.CardService.DeleteCard(userID, id)
}

// CardOwner service delegation methods
func (s *Service) RegisterCardOwner(userID uint, req *models.RegisterOwnerRequest) (*models.CardOwner, error) {
	return s.CardOwnerService.RegisterCardOwner(userID, req)
}

func (s *Service) RegisterMultipleCards(userID uint, req *models.RegisterMultipleCardsRequest) ([]models.CardOwner, error) {
	return s.CardOwnerService.RegisterMultipleCards(userID, req)
}

func (s *Service) GetCardOwnerProfile(userID uint) (*models.CardOwnerWithCard, error) {
	// For backward compatibility, get the first card owner registration
	profiles, err := s.CardOwnerService.GetCardOwnerProfiles(userID)
	if err != nil {
		return nil, err
	}
	if len(profiles) == 0 {
		return nil, errors.New("card owner not found")
	}
	return &profiles[0], nil
}

func (s *Service) GetCardOwnerProfiles(userID uint) ([]models.CardOwnerWithCard, error) {
	return s.CardOwnerService.GetCardOwnerProfiles(userID)
}

func (s *Service) GetAllCardOwners(userID uint) ([]models.CardOwnerWithCard, error) {
	return s.CardOwnerService.GetAllCardOwners(userID)
}

func (s *Service) UpdateCardOwner(userID uint, cardOwnerID uint, req *models.UpdateCardOwnerRequest) (*models.CardOwner, error) {
	return s.CardOwnerService.UpdateCardOwner(userID, cardOwnerID, req)
}

func (s *Service) DeleteCardOwner(userID uint, cardOwnerID uint) error {
	return s.CardOwnerService.DeleteCardOwner(userID, cardOwnerID)
}
