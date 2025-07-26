package repository

import (
	"tiger-fasttrack-card/internal/database"
	"tiger-fasttrack-card/internal/models"
	"errors"
	
	"gorm.io/gorm"
)

type Repository struct {
	DB *database.Database
}

func New(db *database.Database) *Repository {
	return &Repository{
		DB: db,
	}
}

// Example repository methods for cards
// You can add your actual repository methods here when you create models

/*
func (r *Repository) GetAllCards() ([]models.Card, error) {
	var cards []models.Card
	err := r.DB.GetDB().Find(&cards).Error
	return cards, err
}

func (r *Repository) GetCardByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.DB.GetDB().First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *Repository) CreateCard(card *models.Card) error {
	return r.DB.GetDB().Create(card).Error
}

func (r *Repository) UpdateCard(card *models.Card) error {
	return r.DB.GetDB().Save(card).Error
}

func (r *Repository) DeleteCard(id uint) error {
	return r.DB.GetDB().Delete(&models.Card{}, id).Error
}
*/

// User repository methods

func (r *Repository) CreateUser(user *models.User) error {
	return r.DB.GetDB().Create(user).Error
}

func (r *Repository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.GetDB().First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUser(user *models.User) error {
	return r.DB.GetDB().Save(user).Error
}

func (r *Repository) DeleteUser(id uint) error {
	return r.DB.GetDB().Delete(&models.User{}, id).Error
}
