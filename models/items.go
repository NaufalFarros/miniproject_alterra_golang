package models

import (
	"time"

	"gorm.io/gorm"
)

type Items struct {
	gorm.Model
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	Stock       int       `json:"stock" validate:"required"`
	Image       string    `json:"image" validate:"required"`
	CategoryID  int       `json:"category_id" validate:"required"`
	Category    Category  `gorm:"foreignKey:CategoryID" `
	Created_at  time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
