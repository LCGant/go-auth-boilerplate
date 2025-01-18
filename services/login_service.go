package services

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/LCGant/go-auth-boilerplate/models"
)

// Compares a hashed password with the provided plain-text password
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Logs a failed login attempt in the database
func LogFailedLoginAttempt(db *gorm.DB, userID *uint, c *gin.Context) {
	db.Create(&models.FailedLoginAttempt{
		UserID:      userID,
		AttemptTime: time.Now(),
		IPAddress:   c.ClientIP(),
	})
}

// Records a login attempt in the database, whether successful or not
func RecordLoginAttempt(db *gorm.DB, username, ipAddress string, success bool) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	loginHistory := models.LoginHistory{
		UserID:    &user.ID,
		LoginTime: time.Now(),
		IPAddress: ipAddress,
		Success:   success,
	}
	db.Create(&loginHistory)
}
