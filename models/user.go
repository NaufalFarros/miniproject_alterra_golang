package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email      string    `json:"email" validate:"required"`
	Name       string    `json:"name" validate:"required"`
	Password   string    `json:"password" validate:"required,min=6"`
	TableID    int       `json:"table_id" validate:"required"`
	Table      Table     `gorm:"foreignKey:TableID" validate:"-"`
	RoleID     int       `json:"role_id"`
	Role       Roles     `gorm:"foreignKey:RoleID" validate:"-"`
	Created_at time.Time `json:"created_at" gorm:"autoCreateTime"`
	Updated_at time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
