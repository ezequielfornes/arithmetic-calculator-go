package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `gorm:"unique;not null"`
	Password string  `gorm:"not null"`
	Status   string  `gorm:"default:'active';not null"`
	Balance  float64 `gorm:"default:100.0;not null"`
}
