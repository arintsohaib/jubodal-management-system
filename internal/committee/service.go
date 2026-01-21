package committee

import (
	"context"
	"fmt"
	"time"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
)

// Service defines business logic for committees and jurisdictions
type Service struct {
	repo *Repository
}

// NewService creates a new committee service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// JURISDICTIONS

// CreateJurisdiction handles creation of a new jurisdiction with hierarchy checks
func (s *Service) CreateJurisdiction(ctx context.Context, j *models.Jurisdiction) error {
	// 1. If not Central (Level 1), must have a parent
	if j.LevelID > 1 {
		if j.ParentID == nil {
			return fmt.Errorf("jurisdictions at level %d must have a parent", j.LevelID)
		}
		// 2. Validate parent existence and level
		parent, err := s.repo.GetJurisdiction(ctx, *j.ParentID)
		if err != nil {
			return fmt.Errorf("parent jurisdiction not found: %w", err)
		}
		// Parent must be exactly one level above or follow specific rules (e.g. District -> Upazila/Municipality)
		// For simplicity, we check if parent level rank is less than current level rank
		if parent.LevelID >= j.LevelID {
			return fmt.Errorf("parent jurisdiction must be at a higher level (Level %d vs Level %d)", parent.LevelID, j.LevelID)
		}
	} else {
		j.ParentID = nil
	}

	return s.repo.CreateJurisdiction(ctx, j)
}

// ListJurisdictionTree returns jurisdictions under a parent
func (s *Service) ListJurisdictionTree(ctx context.Context, parentID *uuid.UUID) ([]*models.Jurisdiction, error) {
	return s.repo.ListJurisdictions(ctx, nil, parentID)
}

// COMMITTEES

// CreateCommittee handles committee creation logic
func (s *Service) CreateCommittee(ctx context.Context, c *models.Committee) error {
	// 1. Verify jurisdiction exists
	if _, err := s.repo.GetJurisdiction(ctx, c.JurisdictionID); err != nil {
		return fmt.Errorf("jurisdiction not found: %w", err)
	}

	// 2. Check for existing active committee
	active, err := s.repo.GetActiveCommittee(ctx, c.JurisdictionID)
	if err != nil {
		return err
	}
	if active != nil {
		return fmt.Errorf("an active committee already exists for this jurisdiction")
	}

	// 3. Set default status and expiry
	if c.Status == "" {
		c.Status = models.StatusProposed
	}
	
	if c.Type == models.TypeConvener {
		// Convener committees usually last 6 months
		expiry := time.Now().AddDate(0, 6, 0)
		c.ExpiresAt = &expiry
	} else if c.Type == models.TypeFull {
		// Full committees usually last 3 years
		expiry := time.Now().AddDate(3, 0, 0)
		c.ExpiresAt = &expiry
	}

	return s.repo.CreateCommittee(ctx, c)
}

// AddMember adds a member to a committee with size and position constraints
func (s *Service) AddMember(ctx context.Context, m *models.CommitteeMember) error {
	// 1. Get committee and jurisdiction details
	query := `
		SELECT c.id, c.type, c.status, j.level_id 
		FROM committees c 
		JOIN jurisdictions j ON c.jurisdiction_id = j.id 
		WHERE c.id = $1
	`
	var cType, cStatus string
	var levelID int
	err := s.repo.db.QueryRow(ctx, query, m.CommitteeID).Scan(&m.CommitteeID, &cType, &cStatus, &levelID)
	if err != nil {
		return fmt.Errorf("committee not found: %w", err)
	}

	// 2. Load existing members
	members, err := s.repo.GetCommitteeMembers(ctx, m.CommitteeID)
	if err != nil {
		return err
	}

	// 3. Size constraints
	maxSize := 151 // Default for District Full
	if cType == models.TypeConvener {
		maxSize = 11
		if len(members) >= maxSize {
			return fmt.Errorf("convener committee cannot exceed %d members", maxSize)
		}
	} else {
		// Full Committee limits by level (as per docs)
		switch levelID {
		case 3: maxSize = 151 // District
		case 4: maxSize = 101 // Upazila/Municipality
		case 5: maxSize = 71  // Union/Ward
		case 6: maxSize = 31  // Ward
		default: maxSize = 151 // Central/Division
		}
		if len(members) >= maxSize {
			return fmt.Errorf("full committee at level %d cannot exceed %d members", levelID, maxSize)
		}
	}

	// 4. Position uniqueness and duplication
	for _, member := range members {
		if member.UserID == m.UserID {
			return fmt.Errorf("user is already a member of this committee")
		}
		
		// Only 1 President, General Secretary, Convener, Member Secretary
		if m.PositionID == member.PositionID {
			// Check if this position is unique (Rank 1 or 2 usually)
			if m.PositionRank <= 2 {
				return fmt.Errorf("this position is already occupied in this committee")
			}
		}
	}

	return s.repo.AddMember(ctx, m)
}

// IsChildJurisdiction checks if targetID is a sub-unit of parentID (recursive)
func (s *Service) IsChildJurisdiction(ctx context.Context, parentID, targetID uuid.UUID) (bool, error) {
	if parentID == targetID {
		return true, nil
	}

	// Simple recursive check using a CTE
	query := `
		WITH RECURSIVE sub_jurisdictions AS (
			SELECT id, parent_id FROM jurisdictions WHERE id = $1
			UNION ALL
			SELECT j.id, j.parent_id FROM jurisdictions j
			INNER JOIN sub_jurisdictions sj ON j.parent_id = sj.id
		)
		SELECT EXISTS(SELECT 1 FROM sub_jurisdictions WHERE id = $2)
	`
	var exists bool
	err := s.repo.db.QueryRow(ctx, query, parentID, targetID).Scan(&exists)
	return exists, err
}

// ActivateCommittee moves a proposed committee to active and dissolves existing ones
func (s *Service) ActivateCommittee(ctx context.Context, id uuid.UUID, approvedBy uuid.UUID) error {
	// 1. Get the proposed committee
	query := `SELECT id, jurisdiction_id, status FROM committees WHERE id = $1`
	var c models.Committee
	err := s.repo.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.JurisdictionID, &c.Status)
	if err != nil {
		return fmt.Errorf("committee not found: %w", err)
	}

	if c.Status != models.StatusProposed {
		return fmt.Errorf("only proposed committees can be activated")
	}

	// 2. Start transaction
	tx, err := s.repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 3. Dissolve existing active committee in this jurisdiction
	dissolveQuery := `UPDATE committees SET status = 'dissolved', updated_at = NOW() WHERE jurisdiction_id = $1 AND status = 'active'`
	_, err = tx.Exec(ctx, dissolveQuery, c.JurisdictionID)
	if err != nil {
		return fmt.Errorf("failed to dissolve existing committee: %w", err)
	}

	// 4. Activate new committee
	activateQuery := `UPDATE committees SET status = 'active', formed_at = NOW(), approved_by = $1, updated_at = NOW() WHERE id = $2`
	_, err = tx.Exec(ctx, activateQuery, approvedBy, id)
	if err != nil {
		return fmt.Errorf("failed to activate committee: %w", err)
	}

	return tx.Commit(ctx)
}
