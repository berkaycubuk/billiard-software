package models

import (
	"time"
)

type User struct {
	ID              uint              `gorm:"primaryKey" json:"id"`
	Name            string            `gorm:"type:varchar(255)" json:"name"`
	Surname         string            `gorm:"type:varchar(255)" json:"surname"`
	Email           string            `gorm:"index:idx_email,unique,type:varchar(255)" json:"email"`
	Phone           string            `gorm:"index:idx_phone,unique,type:varchar(255)" json:"phone"`
	EmailVerifiedAt *time.Time        `json:"email_verified_at"`
	VerifyToken		string					`gorm:"type:varchar(255)" json:"verify_token"`
	Password        string            `gorm:"type:varchar(255)"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       *time.Time        `json:"deleted_at"`
	Role            UserRole          `json:"role"`
	Subscription    *UserSubscription `json:"subscription"`
}

type PasswordReset struct {
	ID			uint		`gorm:"primaryKey" json:"id"`
	UserID	uint		`json:"user_id"`
	Token		string	`gorm:"type:varchar(255)" json:"token"`
	Valid		bool		`json:"valid"`
	CreatedAt	time.Time	`json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
