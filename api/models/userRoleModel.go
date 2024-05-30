package models

import "time"

type UserRole struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `json:"user_id"`
	RoleID    uint       `json:"role_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Role      Role       `json:"role"`
}
