package models

import (
	"time"

	"gorm.io/gorm"
)

type Table struct {
	gorm.Model
	Table_number string    `json:"table_number" validate:"required"`
	Created_at   time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
