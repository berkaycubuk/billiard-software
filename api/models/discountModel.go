package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	DISCOUNT_TYPE_PRICE uint8 = 1
	DISCOUNT_TYPE_PERCENT	  = 2
)

type Discount struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	Type		uint8		`json:"type"`
	Price		decimal.Decimal `gorm:"type:decimal(10,2)" json:"price"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt	*time.Time	`json:"deleted_at"`
}
