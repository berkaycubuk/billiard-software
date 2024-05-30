package table

import (
	"errors"

	"github.com/berkaycubuk/billiard_software_api/database"
	"github.com/berkaycubuk/billiard_software_api/models"
	"gorm.io/gorm"
)

type Table struct {
	ID uint `json:"id"`
}

type GetAllResponse struct {
	Tables []Table `json:"tables"`
}

type DeleteTableRequest struct {
	ID uint `json:"id" binding:"required" validate:"required"`
}

type UpdateTableRequest struct {
	ID   uint   `json:"id" binding:"required" validate:"required"`
	Name string `json:"name" binding:"required" validate:"required"`
	// Status	uint8	`json:"status" binding:"required" validate:"required"`
}

type NewTableRequest struct {
	Name string `json:"name" binding:"required" validate:"required"`
	// Status uint8  `json:"status" binding:"required" validate:"required"`
}

func create(name string) error {
	table := models.Table{
		Name:   name,
		Status: models.TABLE_ACTIVE,
	}
	err := database.DB.Create(&table).Error
	if err != nil {
		return err
	}

	return nil
}

func find(id uint) (*models.Table, error) {
	var table models.Table
	err := database.DB.Where("id = ?", id).First(&table).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &table, nil
}

func delete(id uint) error {
	err := database.DB.Where("id = ?", id).Delete(&models.Table{}).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}

		return err
	}

	return nil
}
