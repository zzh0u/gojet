package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a student record
type User struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	DeletedAt gorm.DeletedAt `json:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
