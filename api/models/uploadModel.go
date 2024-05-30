package models

import (
	"time"
)

type Upload struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FileName  string    `gorm:"type:varchar(255)" json:"file_name"`
	MimeType  string    `gorm:"type:varchar(255)" json:"mime_type"`
	Size      uint      `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
