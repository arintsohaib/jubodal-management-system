package models

import (
	"time"

	"github.com/google/uuid"
)

// Activity categories
const (
	CategoryPolitical      = "political"
	CategorySocial         = "social"
	CategoryOrganizational = "organizational"
	CategoryProtest        = "protest"
	CategoryOther          = "other"
)

// Activity represents a logged political or organizational activity
type Activity struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	UserID         uuid.UUID  `json:"user_id" db:"user_id"`
	JurisdictionID uuid.UUID  `json:"jurisdiction_id" db:"jurisdiction_id"`
	CommitteeID    *uuid.UUID `json:"committee_id" db:"committee_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Category       string     `json:"category" db:"category"`
	ActivityDate   time.Time  `json:"activity_date" db:"activity_date"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`

	// Joined fields
	UserName         string `json:"user_name,omitempty" db:"user_name"`
	JurisdictionName string `json:"jurisdiction_name,omitempty" db:"jurisdiction_name"`
}

// ActivityProof represents an uploaded file as evidence for an activity
type ActivityProof struct {
	ID         uuid.UUID `json:"id" db:"id"`
	ActivityID uuid.UUID `json:"activity_id" db:"activity_id"`
	FilePath   string    `json:"file_path" db:"file_path"`
	FileType   string    `json:"file_type" db:"file_type"`
	FileSize   int       `json:"file_size" db:"file_size"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// Task statuses
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusVerified   = "verified"
	TaskStatusCancelled  = "cancelled"
)

// Task represents an assigned piece of work
type Task struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	CreatorID      uuid.UUID  `json:"creator_id" db:"creator_id"`
	AssigneeID     *uuid.UUID `json:"assignee_id" db:"assignee_id"`
	CommitteeID    *uuid.UUID `json:"committee_id" db:"committee_id"`
	JurisdictionID uuid.UUID  `json:"jurisdiction_id" db:"jurisdiction_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Status         string     `json:"status" db:"status"`
	Priority       int        `json:"priority" db:"priority"`
	DueDate        *time.Time `json:"due_date" db:"due_date"`
	CompletedAt    *time.Time `json:"completed_at" db:"completed_at"`
	VerifiedAt     *time.Time `json:"verified_at" db:"verified_at"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`
}

// Event represents an organized gathering or program
type Event struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	JurisdictionID uuid.UUID  `json:"jurisdiction_id" db:"jurisdiction_id"`
	OrganizerID    uuid.UUID  `json:"organizer_id" db:"organizer_id"`
	Title          string     `json:"title" db:"title"`
	Description    string     `json:"description" db:"description"`
	Location       string     `json:"location" db:"location"`
	StartTime      time.Time  `json:"start_time" db:"start_time"`
	EndTime        *time.Time `json:"end_time" db:"end_time"`
	IsPublic       bool       `json:"is_public" db:"is_public"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`
}

// EventAttendance tracks participation in an event
type EventAttendance struct {
	ID         uuid.UUID `json:"id" db:"id"`
	EventID    uuid.UUID `json:"event_id" db:"event_id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	AttendedAt time.Time `json:"attended_at" db:"attended_at"`
}
