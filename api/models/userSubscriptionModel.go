package models

import (
	"time"
)

const (
	USER_SUBSCRIPTION_STATUS_ACTIVE   uint8 = 1
	USER_SUBSCRIPTION_STATUS_PAUSED         = 2
	USER_SUBSCRIPTION_STATUS_CANCELED       = 3
	USER_SUBSCRIPTION_STATUS_ENDED          = 4
)

type UserSubscription struct {
	ID             uint       `gorm:"foreignKey" json:"id"`
	UserID         uint       `json:"user_id"`
	SubscriptionID uint       `json:"subscription_id"`
	Name           string     `gorm:"type:varchar(255)" json:"name"`
	Status         uint8      `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	EndingAt       time.Time  `json:"ending_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	User					 *User			`json:"user"`
}

const (
	USER_SUBSCRIPTION_ACTION_ACTIVATED uint8 = 1
	USER_SUBSCRIPTION_ACTION_PAUSED          = 2
	USER_SUBSCRIPTION_ACTION_RENEWED         = 3
)

type UserSubscriptionChunk struct {
	ID                 uint      `gorm:"foreignKey" json:"id"`
	UserSubscriptionID uint      `json:"user_subscription_id"`
	Action             uint8     `json:"action"`
	ActionBy           *uint     `json:"action_by"`
	CreatedAt          time.Time `json:"created_at"`
}
