package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItems struct {
	gorm.Model
	ItemID         uint      `json:"item_id" validate:"required"`
	Item           *Items    `gorm:"foreignKey:ItemID" validate:"-"`
	Quantity       int       `json:"quantity" validate:"required"`
	SubTotal       int       `json:"sub_total" validate:"required"`
	Quantity_total int       `json:"quantity_total" validate:"required"`
	Total_price    int       `json:"total_price" validate:"required"`
	OrdersID       uint      `json:"orders_id" validate:"required"`
	Orders         *Orders   `gorm:"foreignKey:OrdersID" validate:"-"`
	Created_at     time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
