package services

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given plain-text password using bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ValidatePassword checks if the password meets complexity requirements
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return &PasswordValidationError{"Password must be at least 8 characters long."}
	}

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return &PasswordValidationError{"Password must contain at least one uppercase letter."}
	}
	if !hasLower {
		return &PasswordValidationError{"Password must contain at least one lowercase letter."}
	}
	if !hasNumber {
		return &PasswordValidationError{"Password must contain at least one number."}
	}
	if !hasSpecial {
		return &PasswordValidationError{"Password must contain at least one special character."}
	}

	return nil
}

// PasswordValidationError represents an error for invalid password complexity
type PasswordValidationError struct {
	Message string
}

// Error returns the error message for PasswordValidationError
func (e *PasswordValidationError) Error() string {
	return e.Message
}
