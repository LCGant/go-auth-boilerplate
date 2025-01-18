package models

import "time"

// LoginRequest is the model for login requests
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// FailedLoginAttempt is the model for failed login attempts
type FailedLoginAttempt struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      *uint     `gorm:"index"`
	AttemptTime time.Time `gorm:"autoCreateTime"`
	IPAddress   string    `gorm:"size:45"`
}

// LoginHistory is the model for login history
type LoginHistory struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    *uint     `gorm:"index"`
	LoginTime time.Time `gorm:"autoCreateTime"`
	IPAddress string    `gorm:"size:45"`
	Success   bool
}
