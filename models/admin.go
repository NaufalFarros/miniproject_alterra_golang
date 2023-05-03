package models

import (
	"time"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Email      string    `json:"email" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Password   string    `json:"password" validate:"required"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
