package models

import (
	"time"
)

const (
	ORDER_ACTION_INIT uint8 = 1
	ORDER_ACTION_CANCEL     = 2
	ORDER_ACTION_PAY		= 3
	ORDER_ACTION_TRANSFER	= 4
	ORDER_ACTION_APPROVE	= 5
)

const (
	ORDER_METHOD_PHYSICAL uint8 = 1
	ORDER_METHOD_VIPPS	    = 2
	ORDER_METHOD_SYSTEM	    = 3
)

type OrderHistory struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	OrderID		uint		`json:"order_id"`
	Action		uint8		`json:"action"`
	Method		uint8		`json:"method"`
	CreatedAt	time.Time	`json:"created_at"`
	Order		Order
}
