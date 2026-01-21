package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a BJDMS user
type User struct {
	ID                   uuid.UUID  `json:"id" db:"id"`
	FullName             string     `json:"full_name" db:"full_name"`
	FullNameBn           *string    `json:"full_name_bn,omitempty" db:"full_name_bn"`
	NID                  string     `json:"nid" db:"nid"`
	Phone                string     `json:"phone" db:"phone"`
	Email                *string    `json:"email,omitempty" db:"email"`
	PasswordHash         string     `json:"-" db:"password_hash"` // Never expose in JSON
	IsActive             bool       `json:"is_active" db:"is_active"`
	VerifiedAt           *time.Time `json:"verified_at,omitempty" db:"verified_at"`
	FailedLoginAttempts  int        `json:"-" db:"failed_login_attempts"`
	LockedUntil          *time.Time `json:"-" db:"locked_until"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt            *time.Time `json:"-" db:"deleted_at"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Phone    string `json:"phone" validate:"required,phone"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse contains JWT tokens
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds
	User         *User  `json:"user"`
}

// RefreshRequest for token refresh
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuditLog represents a system audit event
type AuditLog struct {
	ID         uuid.UUID              `json:"id" db:"id"`
	UserID     *uuid.UUID             `json:"user_id,omitempty" db:"user_id"`
	Action     string                 `json:"action" db:"action"`
	Entity     string                 `json:"entity" db:"entity"`
	EntityID   *uuid.UUID             `json:"entity_id,omitempty" db:"entity_id"`
	OldValue   map[string]interface{} `json:"old_value,omitempty" db:"old_value"`
	NewValue   map[string]interface{} `json:"new_value,omitempty" db:"new_value"`
	IPAddress  *string                `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent  *string                `json:"user_agent,omitempty" db:"user_agent"`
	Metadata   map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt  time.Time              `json:"created_at" db:"created_at"`
}

// IsLocked checks if user account is locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// IsVerified checks if user has verified status
func (u *User) IsVerified() bool {
	return u.VerifiedAt != nil
}

// MaskPhone returns masked phone number for display
func (u *User) MaskPhone() string {
	if len(u.Phone) < 8 {
		return "***"
	}
	return u.Phone[:6] + "***" + u.Phone[len(u.Phone)-4:]
}
