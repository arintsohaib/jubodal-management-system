package committee

import (
	"context"
	"fmt"
	"time"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for committees and jurisdictions
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new committee repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// JURISDICTIONS

// CreateJurisdiction inserts a new jurisdiction
func (r *Repository) CreateJurisdiction(ctx context.Context, j *models.Jurisdiction) error {
	query := `
		INSERT INTO jurisdictions (level_id, parent_id, name, name_bn, is_urban, population)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		j.LevelID, j.ParentID, j.Name, j.NameBn, j.IsUrban, j.Population,
	).Scan(&j.ID, &j.CreatedAt, &j.UpdatedAt)
}

// GetJurisdiction retrieves a jurisdiction by ID
func (r *Repository) GetJurisdiction(ctx context.Context, id uuid.UUID) (*models.Jurisdiction, error) {
	query := `
		SELECT id, level_id, parent_id, name, name_bn, is_urban, population, created_at, updated_at
		FROM jurisdictions
		WHERE id = $1 AND deleted_at IS NULL
	`
	var j models.Jurisdiction
	err := r.db.QueryRow(ctx, query, id).Scan(
		&j.ID, &j.LevelID, &j.ParentID, &j.Name, &j.NameBn, &j.IsUrban, &j.Population, &j.CreatedAt, &j.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("jurisdiction not found")
	}
	return &j, err
}

// ListJurisdictions returns all jurisdictions at a specific level or parent
func (r *Repository) ListJurisdictions(ctx context.Context, levelID *int, parentID *uuid.UUID) ([]*models.Jurisdiction, error) {
	query := `
		SELECT id, level_id, parent_id, name, name_bn, is_urban, population, created_at, updated_at
		FROM jurisdictions
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	if levelID != nil {
		args = append(args, *levelID)
		query += fmt.Sprintf(" AND level_id = $%d", len(args))
	}
	if parentID != nil {
		args = append(args, *parentID)
		query += fmt.Sprintf(" AND parent_id = $%d", len(args))
	} else if levelID != nil && *levelID > 1 {
		// If level > Central but parent not specified, it's an invalid list request usually
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Jurisdiction
	for rows.Next() {
		var j models.Jurisdiction
		err := rows.Scan(&j.ID, &j.LevelID, &j.ParentID, &j.Name, &j.NameBn, &j.IsUrban, &j.Population, &j.CreatedAt, &j.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, &j)
	}
	return list, nil
}

// COMMITTEES

// CreateCommittee inserts a new committee record
func (r *Repository) CreateCommittee(ctx context.Context, c *models.Committee) error {
	query := `
		INSERT INTO committees (jurisdiction_id, type, status, formed_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		c.JurisdictionID, c.Type, c.Status, c.FormedAt, c.ExpiresAt,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// UpdateCommitteeStatus changes the status of a committee
func (r *Repository) UpdateCommitteeStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `UPDATE committees SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

// GetActiveCommittee retrieves the active committee for a jurisdiction
func (r *Repository) GetActiveCommittee(ctx context.Context, jurisdictionID uuid.UUID) (*models.Committee, error) {
	query := `
		SELECT id, jurisdiction_id, type, status, formed_at, expires_at, created_at, updated_at
		FROM committees
		WHERE jurisdiction_id = $1 AND status = 'active' AND deleted_at IS NULL
	`
	var c models.Committee
	err := r.db.QueryRow(ctx, query, jurisdictionID).Scan(
		&c.ID, &c.JurisdictionID, &c.Type, &c.Status, &c.FormedAt, &c.ExpiresAt, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil // No active committee found
	}
	return &c, err
}

// MEMBERS

// AddMember adds a user to a committee with a position
func (r *Repository) AddMember(ctx context.Context, m *models.CommitteeMember) error {
	query := `
		INSERT INTO committee_members (committee_id, user_id, position_id, joined_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	if m.JoinedAt.IsZero() {
		m.JoinedAt = time.Now()
	}
	return r.db.QueryRow(ctx, query, m.CommitteeID, m.UserID, m.PositionID, m.JoinedAt).Scan(&m.ID, &m.CreatedAt)
}

// GetCommitteeMembers lists all members of a committee
func (r *Repository) GetCommitteeMembers(ctx context.Context, committeeID uuid.UUID) ([]*models.CommitteeMember, error) {
	query := `
		SELECT cm.id, cm.committee_id, cm.user_id, cm.position_id, cm.joined_at, cm.ended_at, cm.is_active,
		       u.full_name as user_name, p.name as position_name, p.rank as position_rank
		FROM committee_members cm
		JOIN users u ON cm.user_id = u.id
		JOIN positions p ON cm.position_id = p.id
		WHERE cm.committee_id = $1 AND cm.ended_at IS NULL
		ORDER BY p.rank ASC
	`
	rows, err := r.db.Query(ctx, query, committeeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.CommitteeMember
	for rows.Next() {
		var m models.CommitteeMember
		err := rows.Scan(
			&m.ID, &m.CommitteeID, &m.UserID, &m.PositionID, &m.JoinedAt, &m.EndedAt, &m.IsActive,
			&m.UserName, &m.PositionName, &m.PositionRank,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, &m)
	}
	return members, nil
}
