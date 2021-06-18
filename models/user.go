package models

import (
	"time"
)

// User model
type User struct {
	ID                uint      `json:"id"`
	Email             *string   `json:"email"`
	Password          *string   `json:"password"`
	MobileNumber      *string   `json:"mobile_number"`
	IsActive          bool      `json:"is_active"`
	IsVerified        bool      `json:"is_verified"`
	HasActivePassword bool      `json:"has_active_password"`
	IsAbsolete        bool      `json:"is_absolete"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TableName gives table name of model
func (u User) TableName() string {
	return "user"
}
