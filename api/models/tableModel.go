package models

import (
	"time"
)

const (
	TABLE_ACTIVE uint8 = 1
	TABLE_INACTIVE     = 2
)

type Table struct {
	ID		uint		`json:"id" gorm:"primaryKey"`
	Name		string		`json:"name"`
	Status		uint8		`json:"status"`
	Games		[]Game		`json:"games"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt	*time.Time	`json:"deleted_at"`
}
