package models

import (
	"time"

	"github.com/google/uuid"
)

// Complaint status constants
const (
	ComplaintStatusReceived    = "received"
	ComplaintStatusUnderReview = "under_review"
	ComplaintStatusActionTaken = "action_taken"
	ComplaintStatusClosed      = "closed"
	ComplaintStatusRejected    = "rejected"
)

// Complaint represents a grievance submitted by a member or the public
type Complaint struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	TrackingID         string     `json:"tracking_id" db:"tracking_id"`
	UserID             *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	JurisdictionID     uuid.UUID  `json:"jurisdiction_id" db:"jurisdiction_id"`
	IsAnonymous        bool       `json:"is_anonymous" db:"is_anonymous"`
	ComplainantName    *string    `json:"complainant_name,omitempty" db:"complainant_name"`
	ComplainantContact *string    `json:"complainant_contact,omitempty" db:"complainant_contact"`
	Subject            string     `json:"subject" db:"subject"`
	Description        string     `json:"description" db:"description"`
	Status             string     `json:"status" db:"status"`
	AnonymousIPHash    *string    `json:"-" db:"anonymous_ip_hash"`
	AssignedToID       *uuid.UUID `json:"assigned_to_id,omitempty" db:"assigned_to_id"`
	ResolutionNotes    *string    `json:"resolution_notes,omitempty" db:"resolution_notes"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt          *time.Time `json:"-" db:"deleted_at"`

	// Joined fields
	JurisdictionName string `json:"jurisdiction_name,omitempty" db:"jurisdiction_name"`
	AssignedToName   string `json:"assigned_to_name,omitempty" db:"assigned_to_name"`
}

// ComplaintEvidence represents a file attached to a complaint
type ComplaintEvidence struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ComplaintID uuid.UUID `json:"complaint_id" db:"complaint_id"`
	FilePath    string    `json:"file_path" db:"file_path"`
	FileType    string    `json:"file_type" db:"file_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// ComplaintLog represents an entry in the complaint's audit trail
type ComplaintLog struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ComplaintID uuid.UUID  `json:"complaint_id" db:"complaint_id"`
	UserID      *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Action      string     `json:"action" db:"action"`
	OldStatus   *string    `json:"old_status,omitempty" db:"old_status"`
	NewStatus   *string    `json:"new_status,omitempty" db:"new_status"`
	Note        *string    `json:"note,omitempty" db:"note"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`

	// Joined fields
	UserName string `json:"user_name,omitempty" db:"user_name"`
}
