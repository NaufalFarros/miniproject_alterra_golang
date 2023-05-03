package models

import (
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	Name_customer string    `json:"name_customer" validate:"required"`
	Phone         string    `json:"phone" validate:"required"`
	Table_number  int       `json:"table_number" validate:"required"`
	ItemsID       int       `json:"items_id" validate:"required"`
	Items         Items     `gorm:"foreignKey:ItemsID" validate:"required"`
	Quantity      int       `json:"quantity" validate:"required"`
	Total_price   int       `json:"total_price" validate:"required"`
	Status_order  string    `json:"status_order" default:"pending" validate:"required"`
	Created_at    time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
