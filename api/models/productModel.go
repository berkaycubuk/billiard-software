package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	ON_SALE      uint8 = 1
	OUT_OF_STOCK       = 2
	DRAFT              = 3
)

type Product struct {
	ID        uint            `gorm:"primaryKey" json:"id"`
	Name      string          `gorm:"type:varchar(255)" json:"name"`
	Price     decimal.Decimal `gorm:"type:decimal(10,2)" json:"price"`
	Status    uint8           `json:"status"`
	Order     int             `json:"order"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at"`
	Image     ProductImage    `gorm:"constraint:OnDelete: CASCADE" json:"image"`
}

type ProductImage struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	ProductID      uint       `json:"product_id"`
	UploadID       uint       `json:"upload_id"`
	UploadFilename string     `gorm:"type:varchar(255)" json:"upload_filename"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
