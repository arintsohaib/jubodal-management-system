package models

import (
	"time"

	"github.com/google/uuid"
)

// JoinRequest statuses
const (
	JoinRequestStatusPending     = "pending"
	JoinRequestStatusUnderReview = "under_review"
	JoinRequestStatusApproved    = "approved"
	JoinRequestStatusRejected    = "rejected"
	JoinRequestStatusCompleted   = "completed"
)

// JoinRequest represents a membership application
type JoinRequest struct {
	ID              uuid.UUID `json:"id" db:"id"`
	FullName        string    `json:"full_name" db:"full_name"`
	FullNameBn      string    `json:"full_name_bn" db:"full_name_bn"`
	Phone           string    `json:"phone" db:"phone"`
	NID             string    `json:"nid" db:"nid"`
	DateOfBirth     time.Time `json:"date_of_birth" db:"date_of_birth"`
	Gender          string    `json:"gender" db:"gender"`
	BloodGroup      string    `json:"blood_group" db:"blood_group"`
	Occupation      string    `json:"occupation" db:"occupation"`
	Address         string    `json:"address" db:"address"`
	JurisdictionID  uuid.UUID `json:"jurisdiction_id" db:"jurisdiction_id"`
	AppliedAt       time.Time `json:"applied_at" db:"applied_at"`
	Status          string    `json:"status" db:"status"`
	ReferredByID    *uuid.UUID `json:"referred_by_id,omitempty" db:"referred_by_id"`
	RejectionReason string    `json:"rejection_reason,omitempty" db:"rejection_reason"`
	ProcessedByID   *uuid.UUID `json:"processed_by_id,omitempty" db:"processed_by_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`

	// Joined fields
	JurisdictionName string `json:"jurisdiction_name,omitempty" db:"jurisdiction_name"`
	ReferrerName     string `json:"referrer_name,omitempty" db:"referrer_name"`
}

// JoinRequestLog tracks the history of an application
type JoinRequestLog struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	RequestID  uuid.UUID  `json:"request_id" db:"request_id"`
	ActorID    uuid.UUID  `json:"actor_id" db:"actor_id"`
	Action     string     `json:"action" db:"action"`
	OldStatus  *string    `json:"old_status,omitempty" db:"old_status"`
	NewStatus  *string    `json:"new_status,omitempty" db:"new_status"`
	Note       string     `json:"note" db:"note"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	
	// Joined fields
	ActorName string `json:"actor_name,omitempty" db:"actor_name"`
}
