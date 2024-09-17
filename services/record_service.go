package services

import (
	"arithmetic-calculator/config"
	"arithmetic-calculator/models"
	"errors"
	"gorm.io/gorm"
	"math"
)

type PaginatedRecords struct {
	Records    []models.Record `json:"records"`
	Page       int             `json:"page"`
	TotalPages int             `json:"total_pages"`
	TotalItems int64           `json:"total_items"`
}

func FetchRecords(username string, page, pageSize int) (PaginatedRecords, error) {
	var user models.User
	var records []models.Record
	var totalItems int64

	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return PaginatedRecords{}, errors.New("user not found")
		}
		return PaginatedRecords{}, err
	}

	offset := (page - 1) * pageSize

	if err := config.DB.Model(&models.Record{}).Where("user_id = ?", user.ID).Count(&totalItems).Error; err != nil {
		return PaginatedRecords{}, err
	}

	if err := config.DB.Where("user_id = ?", user.ID).Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return PaginatedRecords{}, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return PaginatedRecords{
		Records:    records,
		Page:       page,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}, nil
}
