package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/LCGant/go-auth-boilerplate/models"
	"github.com/LCGant/go-auth-boilerplate/services"
)

// LoginHandler handles the login request
func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		var user models.User
		if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
			services.LogFailedLoginAttempt(db, nil, c)
			services.RecordLoginAttempt(db, req.Username, c.ClientIP(), false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		var userEmail models.UserEmail
		if err := db.Where("user_id = ? AND current = true", user.ID).First(&userEmail).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		if !userEmail.Verified {
			if time.Now().After(userEmail.ExpiresAt) {
				newToken, err := services.GenerateValidateToken()
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating new token"})
					return
				}
				userEmail.VerificationToken = newToken
				userEmail.ExpiresAt = time.Now().Add(24 * time.Hour)
				if err := db.Save(&userEmail).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating token"})
					return
				}
				if err := services.SendEmailVerification(db, user.Email); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending verification email"})
					return
				}
				c.JSON(http.StatusForbidden, gin.H{"message": "Email not verified. A new verification email has been sent."})
				services.RecordLoginAttempt(db, req.Username, c.ClientIP(), false)
				return
			}
			c.JSON(http.StatusForbidden, gin.H{"error": "Email not verified"})
			services.RecordLoginAttempt(db, req.Username, c.ClientIP(), false)
			return
		}

		if !services.ComparePassword(user.PasswordHash, req.Password) {
			services.LogFailedLoginAttempt(db, &user.ID, c)
			services.RecordLoginAttempt(db, req.Username, c.ClientIP(), false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		userIDStr := strconv.FormatUint(uint64(user.ID), 10)
		token, err := services.GenerateToken(userIDStr)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		// Set the cookie
		// TODO: In production, ensure the following:
		// - Use "Secure" cookies when serving over HTTPS (`Secure: true` in SetCookie).
		// - Replace "localhost" with the domain of your application.
		// - Adjust SameSite policy as needed (e.g., "None" for cross-site cookies with HTTPS).
		// - Consider encrypting the cookie value for added security.
		c.SetCookie("token", token, 3600, "/", "localhost", false, true)
		c.Header("Set-Cookie", "token="+token+"; Path=/; Max-Age=3600; HttpOnly; SameSite=Lax")

		services.RecordLoginAttempt(db, req.Username, c.ClientIP(), true)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	}
}
