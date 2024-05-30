package models

import (
	"github.com/shopspring/decimal"
)

type OrderDiscount struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	OrderID		uint		`json:"order_id"`
	DiscountID	*uint		`json:"discount_id"`
	Type		uint8		`json:"type"`
	DiscountPrice	decimal.Decimal `gorm:"type:decimal(10,2)" json:"discount_price"`
	Description	*string		`gorm:"type:varchar(255)" json:"description"`
	Order		Order
	Discount	*Discount
}
