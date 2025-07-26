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

// Repository methods for cards
func (r *Repository) GetAllCards() ([]models.Card, error) {
	var cards []models.Card
	err := r.DB.GetDB().Find(&cards).Error
	return cards, err
}

func (r *Repository) GetCardByID(id uint) (*models.Card, error) {
	var card models.Card
	err := r.DB.GetDB().First(&card, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("card not found")
		}
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

func (r *Repository) GetCardByName(cardName string) (*models.Card, error) {
	var card models.Card
	err := r.DB.GetDB().Where("card_name = ?", cardName).First(&card).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("card not found")
		}
		return nil, err
	}
	return &card, nil
}

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

// CardOwner repository methods

func (r *Repository) CreateCardOwner(owner *models.CardOwner) error {
	return r.DB.GetDB().Create(owner).Error
}

func (r *Repository) GetCardOwnerByID(id uint) (*models.CardOwner, error) {
	var owner models.CardOwner
	err := r.DB.GetDB().Preload("User").Preload("Card").First(&owner, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("card owner not found")
		}
		return nil, err
	}
	return &owner, nil
}

func (r *Repository) GetCardOwnerByUserID(userID uint) (*models.CardOwner, error) {
	var owner models.CardOwner
	err := r.DB.GetDB().Preload("User").Preload("Card").Where("user_id = ?", userID).First(&owner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("card owner not found")
		}
		return nil, err
	}
	return &owner, nil
}

func (r *Repository) GetCardOwnersByUserID(userID uint) ([]models.CardOwner, error) {
	var owners []models.CardOwner
	err := r.DB.GetDB().Preload("User").Preload("Card").Where("user_id = ?", userID).Find(&owners).Error
	if err != nil {
		return nil, err
	}
	return owners, nil
}

func (r *Repository) GetCardOwnerByCardNumberAndCardID(cardNumber string, cardID uint) (*models.CardOwner, error) {
	var owner models.CardOwner
	err := r.DB.GetDB().Preload("User").Preload("Card").Where("card_number = ? AND card_id = ?", cardNumber, cardID).First(&owner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("card owner not found")
		}
		return nil, err
	}
	return &owner, nil
}

func (r *Repository) GetCardOwnerByIDCard(idCard string) (*models.CardOwner, error) {
	var owner models.CardOwner
	err := r.DB.GetDB().Preload("User").Preload("Card").Where("id_card = ?", idCard).First(&owner).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("card owner not found")
		}
		return nil, err
	}
	return &owner, nil
}

func (r *Repository) GetAllCardOwners() ([]models.CardOwner, error) {
	var owners []models.CardOwner
	err := r.DB.GetDB().Preload("User").Preload("Card").Find(&owners).Error
	return owners, err
}

func (r *Repository) UpdateCardOwner(owner *models.CardOwner) error {
	return r.DB.GetDB().Save(owner).Error
}

func (r *Repository) DeleteCardOwner(id uint) error {
	return r.DB.GetDB().Delete(&models.CardOwner{}, id).Error
}
