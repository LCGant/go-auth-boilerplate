package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/LCGant/go-auth-boilerplate/models"
	"github.com/LCGant/go-auth-boilerplate/services"
)

// ResetPasswordHandler resets the password of a user
func ResetPasswordHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Token       string `json:"token" binding:"required"`
			NewPassword string `json:"newPassword" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
			return
		}

		userID, err := services.ValidatePasswordResetToken(db, request.Token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		newPasswordHash, err := services.HashPassword(request.NewPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing new password", "details": err.Error()})
			return
		}

		user.PasswordHash = newPasswordHash
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating password", "details": err.Error()})
			return
		}

		// Remove o token usado
		if err := db.Exec("DELETE FROM password_reset_tokens WHERE token = ?", request.Token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting token", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
	}
}

// VerifyTokenHandler verifies the password reset token
func VerifyTokenHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.Redirect(http.StatusFound, "/error/401")
			return
		}
		_, err := services.ValidatePasswordResetToken(db, token)
		if err != nil {
			c.Redirect(http.StatusFound, "/error/401")
			return
		}
		c.File("public/pages/reset-password.html")
	}
}

// VerifyEmailTokenHandler verifies the email verification token
func VerifyEmailHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			return
		}

		var emailChangeToken models.EmailChangeToken
		err := db.Where("token = ? AND expires_at > NOW() AND used = FALSE", token).First(&emailChangeToken).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if err := db.Model(&emailChangeToken).Update("used", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating token status"})
			return
		}

		if err := db.Model(&models.User{}).Where("id = ?", emailChangeToken.UserID).Update("email", emailChangeToken.NewEmail).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating user email"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email verified and updated successfully"})
	}
}

// VerifyEmailTokenHandler verifies the email verification token
func VerifyEmailTokenHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
			return
		}

		var userEmail models.UserEmail
		err := db.Where("verification_token = ? AND expires_at > NOW() AND verified = 0", token).First(&userEmail).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if err := db.Model(&userEmail).Update("verified", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating verification status"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
	}
}
