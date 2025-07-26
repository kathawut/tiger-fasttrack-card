package models

import (
	"time"

	"gorm.io/gorm"
)

// Card represents a card in the system (Master Data)
type Card struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	CardName     string         `json:"card_name" gorm:"not null;uniqueIndex"` // Master data - unique card names
	CardImage    string         `json:"card_image" gorm:"not null"`
	CardQuantity int            `json:"card_quantity" gorm:"not null;default:0"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// CreateCardRequest represents the request body for creating a card
type CreateCardRequest struct {
	CardName     string `json:"card_name" binding:"required"`
	CardImage    string `json:"card_image" binding:"required"`
	CardQuantity int    `json:"card_quantity"`
}

// UpdateCardRequest represents the request body for updating a card
type UpdateCardRequest struct {
	CardName     string `json:"card_name"`
	CardImage    string `json:"card_image"`
	CardQuantity *int   `json:"card_quantity,omitempty" binding:"omitempty,min=0"`
}
