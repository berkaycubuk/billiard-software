package models

import "time"

type Role struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	Name		string		`gorm:"type:varchar(255)" json:"name"`
	CreatedAt	time.Time	`json:"created_at"`
	UpdatedAt	time.Time	`json:"updated_at"`
	DeletedAt	*time.Time	`json:"deleted_at"`
	RoleConfigs	*[]RoleConfig	`json:"role_configs"`
}

type RoleConfig struct {
	ID		uint		`gorm:"primaryKey" json:"id"`
	// convert this into RoleID
	RoleId		uint		`json:"role_id"`
	Name		string		`gorm:"type:varchar(255)" json:"name"`
	Value		string		`gorm:"type:varchar(255)" json:"value"`
	ValueType	string		`gorm:"type:varchar(255)" json:"value_type"`
}
