package service

import (
	"errors"
	"tiger-fasttrack-card/internal/models"
	"tiger-fasttrack-card/internal/repository"
	"tiger-fasttrack-card/internal/utils"
)

type Service struct {
	Repo       *repository.Repository
	JWTManager *utils.JWTManager
}

func New(repo *repository.Repository, jwtSecret string) *Service {
	return &Service{
		Repo:       repo,
		JWTManager: utils.NewJWTManager(jwtSecret),
	}
}

// Example service methods
// Add your business logic here when you create models

/*
func (s *Service) GetAllCards() ([]models.Card, error) {
	return s.Repo.GetAllCards()
}

func (s *Service) GetCardByID(id uint) (*models.Card, error) {
	return s.Repo.GetCardByID(id)
}

func (s *Service) CreateCard(card *models.Card) error {
	// Add validation logic here
	return s.Repo.CreateCard(card)
}

func (s *Service) UpdateCard(id uint, card *models.Card) error {
	// Add validation logic here
	card.ID = id
	return s.Repo.UpdateCard(card)
}

func (s *Service) DeleteCard(id uint) error {
	return s.Repo.DeleteCard(id)
}
*/

// Authentication services

func (s *Service) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if username is taken
	existingUser, _ := s.Repo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username is already taken")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		Role:      "user",
	}

	err = s.Repo.CreateUser(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

func (s *Service) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username
	user, err := s.Repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid username or password")
	}

	// Generate tokens
	token, err := s.JWTManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refreshToken, err := s.JWTManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.JWTManager.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Generate new access token
	newToken, err := s.JWTManager.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", errors.New("failed to generate new token")
	}

	return newToken, nil
}

func (s *Service) GetUserProfile(userID uint) (*models.User, error) {
	return s.Repo.GetUserByID(userID)
}

func (s *Service) UpdateUserProfile(userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if username is taken by another user
	if req.Username != "" && req.Username != user.Username {
		existingUser, _ := s.Repo.GetUserByUsername(req.Username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("username is already taken")
		}
		user.Username = req.Username
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}

	err = s.Repo.UpdateUser(user)
	if err != nil {
		return nil, errors.New("failed to update profile")
	}

	return user, nil
}

func (s *Service) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	user, err := s.Repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !utils.CheckPassword(req.CurrentPassword, user.Password) {
		return errors.New("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	user.Password = hashedPassword
	err = s.Repo.UpdateUser(user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}
