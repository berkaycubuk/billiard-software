package models

import (
	"time"
)

const (
	NOTIFICATION_MEDIUM_MAIL uint8 = 1
	NOTIFICATION_MEDIUM_APP = 2
)

const (
	NOTIFICATION_TYPE_SIMPLE uint8 = 1
	NOTIFICATION_TYPE_JOIN_GAME    = 2
)

type Notification struct {
	ID			uint		`gorm:"primaryKey" json:"id"`
	UserID	uint		`json:"user_id"`
	Medium	uint8		`json:"medium"`
	Type		uint8		`json:"type"`
	Message	string	`json:"message"`
	CreatedAt	time.Time	`json:"created_at"`
	ReadAt	*time.Time	`json:"read_at"`
}
