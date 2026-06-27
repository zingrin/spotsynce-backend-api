package models

import (
	"time"

	"gorm.io/gorm"
)

const (
	RoleDriver = "driver"
	RoleAdmin  = "admin"
)

// User represents the users table.
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Role      string         `gorm:"size:20;not null;default:driver;index" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}
