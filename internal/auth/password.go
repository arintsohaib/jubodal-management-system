package auth

import (
	"errors"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost       = 12
	minPasswordLength = 8
)

var (
	ErrWeakPassword    = errors.New("password must be at least 8 characters and contain uppercase, lowercase, and digit")
	ErrInvalidPhone    = errors.New("invalid Bangladesh phone number format")
	phoneRegex         = regexp.MustCompile(`^\+880\d{10}$`)
)

// HashPassword generates bcrypt hash from password
func HashPassword(password string) (string, error) {
	if err := ValidatePassword(password); err != nil {
		return "", err
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// ComparePassword verifies password against hash
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ValidatePassword checks password strength
func ValidatePassword(password string) error {
	if len(password) < minPasswordLength {
		return ErrWeakPassword
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber {
		return ErrWeakPassword
	}

	return nil
}

// ValidatePhone checks Bangladesh phone format
func ValidatePhone(phone string) error {
	if !phoneRegex.MatchString(phone) {
		return ErrInvalidPhone
	}
	return nil
}

// MaskPhone returns masked phone number for display
func MaskPhone(phone string) string {
	if len(phone) < 8 {
		return "***"
	}
	return phone[:6] + "***" + phone[len(phone)-4:]
}
