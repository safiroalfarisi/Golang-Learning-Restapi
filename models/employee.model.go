package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Position  string         `gorm:"type:varchar(100);not null" json:"position"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // Password disembunyikan dari JSON
	Role      string         `gorm:"type:varchar(20);default:staff" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}