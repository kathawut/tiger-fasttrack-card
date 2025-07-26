package models

import (
	"time"

	"gorm.io/gorm"
)

// CardOwner represents a card owner in the system
type CardOwner struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CardID      uint           `json:"card_id" gorm:"not null;index"` // References Card master data by ID
	Card        Card           `json:"card" gorm:"foreignKey:CardID"`
	CardNumber  string         `json:"card_number" gorm:"not null;uniqueIndex:idx_card_number_card_id"`
	IDCard      string         `json:"id_card" gorm:"not null;index"` // Removed uniqueIndex to allow same person to have multiple cards
	PhoneNumber string         `json:"phone_number" gorm:"not null"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// RegisterOwnerRequest represents the request body for registering a card owner
type RegisterOwnerRequest struct {
	CardID      uint   `json:"card_id" binding:"required"`
	CardNumber  string `json:"card_number" binding:"required"`
	IDCard      string `json:"id_card" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// CardRegistration represents a single card registration item
type CardRegistration struct {
	CardID     uint   `json:"card_id" binding:"required"`
	CardNumber string `json:"card_number" binding:"required"`
}

// RegisterMultipleCardsRequest represents the request body for registering multiple cards
type RegisterMultipleCardsRequest struct {
	Cards       []CardRegistration `json:"cards" binding:"required,min=1"`
	IDCard      string             `json:"id_card" binding:"required"`
	PhoneNumber string             `json:"phone_number" binding:"required"`
}

// UpdateCardOwnerRequest represents the request body for updating a card owner
type UpdateCardOwnerRequest struct {
	CardID      uint   `json:"card_id"`
	CardNumber  string `json:"card_number"`
	IDCard      string `json:"id_card"`
	PhoneNumber string `json:"phone_number"`
}

// CardOwnerResponse represents the response for card owner operations
type CardOwnerResponse struct {
	ID          uint      `json:"id"`
	CardID      uint      `json:"card_id"`
	CardNumber  string    `json:"card_number"`
	IDCard      string    `json:"id_card"`
	PhoneNumber string    `json:"phone_number"`
	UserID      uint      `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CardOwnerWithCard represents a card owner with the associated card master data
type CardOwnerWithCard struct {
	CardOwner
	Card *Card `json:"card"`
}
