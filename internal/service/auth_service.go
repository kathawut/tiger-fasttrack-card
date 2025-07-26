package service

import (
	"errors"
	"tiger-fasttrack-card/internal/models"
	"tiger-fasttrack-card/internal/repository"
	"tiger-fasttrack-card/internal/utils"
)

// AuthService handles authentication and user management operations
type AuthService struct {
	repo       *repository.Repository
	jwtManager *utils.JWTManager
}

// NewAuthService creates a new AuthService instance
func NewAuthService(repo *repository.Repository, jwtManager *utils.JWTManager) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user account
func (s *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if username is taken
	existingUser, _ := s.repo.GetUserByUsername(req.Username)
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

	err = s.repo.CreateUser(user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username
	user, err := s.repo.GetUserByUsername(req.Username)
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
	token, err := s.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	return &models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         *user,
	}, nil
}

// RefreshToken generates a new access token from a refresh token
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}

	// Generate new access token
	newToken, err := s.jwtManager.GenerateToken(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		return "", errors.New("failed to generate new token")
	}

	return newToken, nil
}

// GetUserProfile retrieves user profile information
func (s *AuthService) GetUserProfile(userID uint) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

// UpdateUserProfile updates user profile information
func (s *AuthService) UpdateUserProfile(userID uint, req *models.UpdateProfileRequest) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Check if username is taken by another user
	if req.Username != "" && req.Username != user.Username {
		existingUser, _ := s.repo.GetUserByUsername(req.Username)
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

	err = s.repo.UpdateUser(user)
	if err != nil {
		return nil, errors.New("failed to update profile")
	}

	return user, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID uint, req *models.ChangePasswordRequest) error {
	user, err := s.repo.GetUserByID(userID)
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
	err = s.repo.UpdateUser(user)
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

// ValidateUserAccess checks if user exists and is active
func (s *AuthService) ValidateUserAccess(userID uint) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}
	return user, nil
}
