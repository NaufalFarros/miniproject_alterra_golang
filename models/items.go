package models

import (
	"gorm.io/gorm"
)

type Items struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Stock       int    `json:"stock"`
	Image       string `json:"image"`
	CategoryID  int    `json:"category_id"`
	Category Category `gorm:"foreignKey:CategoryID"`
}
