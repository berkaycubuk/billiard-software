package models

import (
	"github.com/shopspring/decimal"
)

const (
	PRODUCT_TYPE_SHOP uint8		= 1
	PRODUCT_TYPE_GAME		= 2
	PRODUCT_TYPE_SUBSCRIPTION	= 3
)

type OrderItem struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	OrderID		uint		` json:"order_id"`
	ProductType	uint8		`json:"product_type"`
	ProductID	*uint		`json:"product_id"`
	ProductName	string		`gorm:"type:varchar(255)" json:"product_name"`
	ProductAmount	int		`json:"product_amount"`
	ProductPrice	decimal.Decimal `gorm:"type:decimal(10,2)" json:"product_price"`
}
