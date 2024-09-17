package models

import "gorm.io/gorm"

type Operation struct {
	gorm.Model
	Type string  `gorm:"unique;not null"`
	Cost float64 `gorm:"not null"`
}
