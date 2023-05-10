package models

import (
	"time"

	"gorm.io/gorm"
)

type Orders struct {
	gorm.Model
	Name_customer string    `json:"name_customer" validate:"required"`
	Phone         string    `json:"phone" validate:"required"`
	Table_number  string    `json:"table_number" validate:"required"`
	Status_order  string    `json:"status_order" default:"pending" validate:"required"`
	UserID        int       `json:"user_id" validate:"required"`
	User          User      `gorm:"foreignKey:UserID" validate:"-"`
	Created_at    time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
