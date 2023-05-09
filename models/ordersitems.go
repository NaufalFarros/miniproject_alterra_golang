package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderItems struct {
	gorm.Model
	ItemID         int       `json:"item_id"`
	Item           Items     `gorm:"foreignKey:ItemID"`
	UserID         int       `json:"user_id"`
	User           User      `gorm:"foreignKey:UserID"`
	Quantity       int       `json:"quantity"`
	SubTotal       int       `json:"sub_total"`
	Quantity_total int       `json:"quantity_total"`
	Total_price    int       `json:"total_price"`
	OrdersID       int       `json:"orders_id"`
	Orders         Orders    `gorm:"foreignKey:OrdersID"`
	Created_at     time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
