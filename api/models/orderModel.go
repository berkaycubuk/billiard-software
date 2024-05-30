package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	ORDER_STATUS_WAITING uint8 = 1
	ORDER_STATUS_CANCELED      = 2
	ORDER_STATUS_PAID          = 3
	ORDER_STATUS_DELETED       = 4
	ORDER_STATUS_PAY_LATER     = 5
	ORDER_STATUS_TRANSFERRED   = 6
	ORDER_STATUS_APPROVED      = 7
	ORDER_STATUS_CAPTURED      = 8
)

type Order struct {
	ID		uint			`gorm:"primaryKey" json:"id"`
	Reference	string			`gorm:"type:varchar(255)" json:"reference"`
	Status		uint8			`json:"status"`
	UserID		*uint			`json:"user_id"`
	Price		decimal.Decimal		`gorm:"type:decimal(10,2)" json:"price"`
	CreatedAt	time.Time		`json:"created_at"`
	UpdatedAt	time.Time		`json:"updated_at"`
	DeletedAt	*time.Time		`json:"deleted_at"`
	Items		[]OrderItem		`gorm:"constraint:OnDelete: CASCADE" json:"items"`
	Detail		OrderDetail		`gorm:"constraint:OnDelete: CASCADE" json:"detail"`
	Discounts	[]OrderDiscount		`gorm:"constraint:OnDelete: CASCADE" json:"discounts"`
	Histories	[]OrderHistory		`gorm:"constraint:OnDelete: CASCADE" json:"histories"`
}
