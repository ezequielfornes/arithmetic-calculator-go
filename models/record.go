package models

import (
	"gorm.io/gorm"
	"time"
)

type Record struct {
	gorm.Model
	UserID            uint
	OperationID       uint
	OperationType     string
	Amount            float64
	UserBalance       float64
	OperationResponse string
	Date              *time.Time
}
