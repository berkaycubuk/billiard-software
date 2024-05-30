package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Pricing struct {
	ID             uint            `gorm:"primaryKey" json:"id"`
	RoleID         *uint           `json:"role_id"`
	SubscriptionID *uint           `json:"subscription_id"`
	PlayerCount    int8            `json:"player_count"`
	PerMinute      decimal.Decimal `gorm:"type:decimal(10,2)" json:"per_minute"`
	Priority       int8            `json:"priority"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *time.Time      `json:"deleted_at"`
	Subscription   *Subscription   `json:"subscription"`
	Role           *Role           `json:"role"`
}
