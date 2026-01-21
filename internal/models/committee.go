package models

import (
	"time"

	"github.com/google/uuid"
)

// JurisdictionLevel represents an administrative level
type JurisdictionLevel struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Rank        int    `json:"rank" db:"rank"`
	Description string `json:"description" db:"description"`
}

// Jurisdiction represents a geographical or administrative unit
type Jurisdiction struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	LevelID     int        `json:"level_id" db:"level_id"`
	ParentID    *uuid.UUID `json:"parent_id" db:"parent_id"`
	Name        string     `json:"name" db:"name"`
	NameBn      *string    `json:"name_bn" db:"name_bn"`
	IsUrban     bool       `json:"is_urban" db:"is_urban"`
	Population  int        `json:"population" db:"population"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
}

// Position represents a role in a committee
type Position struct {
	ID            int    `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	NameBn        string `json:"name_bn" db:"name_bn"`
	Rank          int    `json:"rank" db:"rank"`
	CommitteeType string `json:"committee_type" db:"committee_type"`
	Description   string `json:"description" db:"description"`
}

// Committee status and type constants
const (
	StatusProposed  = "proposed"
	StatusActive    = "active"
	StatusDissolved = "dissolved"
	StatusExpired   = "expired"

	TypeFull     = "full"
	TypeConvener = "convener"
)

// Committee represents a Jubodal committee
type Committee struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	JurisdictionID uuid.UUID  `json:"jurisdiction_id" db:"jurisdiction_id"`
	Type           string     `json:"type" db:"type"`
	Status         string     `json:"status" db:"status"`
	FormedAt       *time.Time `json:"formed_at" db:"formed_at"`
	ExpiresAt      *time.Time `json:"expires_at" db:"expires_at"`
	ApprovedBy     *uuid.UUID `json:"approved_by" db:"approved_by"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt      *time.Time `json:"-" db:"deleted_at"`
}

// CommitteeMember represents a user assigned to a committee
type CommitteeMember struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CommitteeID uuid.UUID  `json:"committee_id" db:"committee_id"`
	UserID      uuid.UUID  `json:"user_id" db:"user_id"`
	PositionID  int        `json:"position_id" db:"position_id"`
	JoinedAt    time.Time  `json:"joined_at" db:"joined_at"`
	EndedAt     *time.Time `json:"ended_at" db:"ended_at"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`

	// Joined data
	UserName     string `json:"user_name,omitempty" db:"user_name"`
	PositionName string `json:"position_name,omitempty" db:"position_name"`
	PositionRank int    `json:"position_rank,omitempty" db:"position_rank"`
}
