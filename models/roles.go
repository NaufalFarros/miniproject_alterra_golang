package models

import (
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	gorm.Model
	Name       string    `json:"name" validate:"required"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
