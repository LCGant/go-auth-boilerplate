package models

import "time"

// User struct
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"unique;not null"`
	Email        string    `gorm:"unique;not null"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	FullName     string
	MobileNumber string
	BirthDate    time.Time
	Gender       string
	ProfileImg   []byte `gorm:"type:MEDIUMBLOB;default:NULL"`
}

// UserEmail struct
type UserEmail struct {
	ID                uint   `gorm:"primaryKey"`
	UserID            *uint  `gorm:"index"`
	Email             string `gorm:"not null"`
	Current           bool
	Verified          bool
	VerificationToken string
	ExpiresAt         time.Time `gorm:"autoUpdateTime"`
}

// EmailChangeToken struct
type EmailChangeToken struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"`
	NewEmail  string `gorm:"not null"`
	Token     string `gorm:"unique;not null"`
	ExpiresAt time.Time
	Used      bool `gorm:"default:false"`
}
