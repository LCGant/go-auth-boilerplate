package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"time"

	"gorm.io/gorm"

	"github.com/LCGant/go-auth-boilerplate/models"
)

// Sends the verification email (already existing)
func SendEmailVerification(db *gorm.DB, email string) error {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	var expiredEmails []models.UserEmail
	if err := db.Where("user_id = ? AND expires_at < ?", user.ID, time.Now()).Find(&expiredEmails).Error; err != nil {
		return err
	}
	if len(expiredEmails) > 3 {
		sort.Slice(expiredEmails, func(i, j int) bool {
			return expiredEmails[i].ExpiresAt.Before(expiredEmails[j].ExpiresAt)
		})
		db.Delete(&expiredEmails[0])
	}

	var existingToken models.UserEmail
	if err := db.Where("user_id = ? AND expires_at >= ?", user.ID, time.Now()).First(&existingToken).Error; err == nil {
		db.Delete(&existingToken)
	}

	token, err := GenerateValidateToken()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(24 * time.Hour)

	if err := db.Model(&models.UserEmail{}).
		Where("user_id = ? AND current = TRUE", user.ID).
		Update("current", false).Error; err != nil {
		return err
	}

	userEmail := models.UserEmail{
		UserID:            &user.ID,
		Email:             email,
		Current:           true,
		Verified:          false,
		VerificationToken: token,
		ExpiresAt:         expiresAt,
	}
	if err := db.Create(&userEmail).Error; err != nil {
		return err
	}

	// TODO: Replace "http://localhost:8080" with a value from environment variables for production.
	verificationLink := fmt.Sprintf("%s/verify-email-token?token=%s", getBaseURL(), token)
	subject := "Email Verification"
	message := fmt.Sprintf("Click the following link to verify your email:\n%s", verificationLink)

	return SendEmailToPythonServer(email, subject, message)
}

// Sends the password reset email (Forgot Password)
func SendPasswordResetEmailByEmail(db *gorm.DB, email string) error {
	var user models.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil
	}

	token, err := GenerateValidateToken()
	if err != nil {
		return err
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	err = db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token, expires_at)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			token = VALUES(token),
			expires_at = VALUES(expires_at)
	`, user.ID, token, expiresAt).Error
	if err != nil {
		return err
	}

	// TODO: Replace "http://localhost:8080" with a value from environment variables for production.
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", getBaseURL(), token)

	subject := "Password Reset Request"
	message := fmt.Sprintf("Click the following link to reset your password: %s", resetLink)

	return SendEmailToPythonServer(email, subject, message)
}

// Sends JSON to the Python server at localhost:5000
func SendEmailToPythonServer(to, subject, message string) error {
	emailData := map[string]string{
		"email":   to,
		"subject": subject,
		"message": message,
	}
	jsonData, err := json.Marshal(emailData)
	if err != nil {
		return err
	}

	// TODO: Replace "http://localhost:5000" with a value from environment variables for production.
	pythonServerURL := os.Getenv("PYTHON_SERVER_URL")
	if pythonServerURL == "" {
		pythonServerURL = "http://localhost:5000"
	}

	resp, err := http.Post(fmt.Sprintf("%s/send-email", pythonServerURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// Validates the reset password token (already existing)
func ValidatePasswordResetToken(db *gorm.DB, token string) (int, error) {
	var result struct {
		UserID int
	}
	err := db.Raw(`
		SELECT user_id 
		FROM password_reset_tokens 
		WHERE token = ? 
		  AND expires_at > NOW()
	`, token).Scan(&result).Error

	if err != nil || result.UserID == 0 {
		return 0, fmt.Errorf("invalid or expired token")
	}

	return result.UserID, nil
}

// Helper function to get the base URL
func getBaseURL() string {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080" // Default to localhost for development
	}
	return baseURL
}
