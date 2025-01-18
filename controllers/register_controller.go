package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/LCGant/go-auth-boilerplate/models"
	"github.com/LCGant/go-auth-boilerplate/services"
)

// RegisterHandler handles the register request
func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		var existingUserByUsername models.User
		if err := db.Where("username = ?", req.Username).First(&existingUserByUsername).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}

		var existingUserByEmail models.User
		if err := db.Where("email = ?", req.Email).First(&existingUserByEmail).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
			return
		}

		birthDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format"})
			return
		}
		age := time.Now().Year() - birthDate.Year()
		if age < 10 || (age == 10 && time.Now().Before(birthDate.AddDate(10, 0, 0))) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You must be at least 10 years old to register"})
			return
		}

		if err := services.ValidatePassword(req.Password); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		passwordHash, err := services.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}

		user := models.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: passwordHash,
			FullName:     req.FullName,
			MobileNumber: req.MobileNumber,
			BirthDate:    birthDate,
			Gender:       req.Gender,
		}

		if err := services.SendEmailVerification(db, user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email. Please try again later."})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Registration successful. Please check your email to verify your account."})
	}
}
