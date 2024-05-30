package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Subscription struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	Name		string		`gorm:"type:varchar(255)" json:"name"`
	Price		decimal.Decimal	`gorm:"type:decimal(10,2)" json:"price"`
	Hours		int		`json:"hours"`
	Role		uint		`json:"role"`
	Hidden		bool	`json:"hidden"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt	*time.Time	`json:"deleted_at"`
}
