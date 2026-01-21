package join

import (
	"context"
	"fmt"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for join requests
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new join repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// Create inserts a new join request
func (r *Repository) Create(ctx context.Context, jr *models.JoinRequest) error {
	query := `
		INSERT INTO join_requests (
			full_name, full_name_bn, phone, nid, date_of_birth, gender, blood_group, 
			occupation, address, jurisdiction_id, referred_by_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, applied_at, status, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		jr.FullName, jr.FullNameBn, jr.Phone, jr.NID, jr.DateOfBirth, jr.Gender, jr.BloodGroup,
		jr.Occupation, jr.Address, jr.JurisdictionID, jr.ReferredByID,
	).Scan(&jr.ID, &jr.AppliedAt, &jr.Status, &jr.CreatedAt, &jr.UpdatedAt)
}

// GetByID retrieves a single join request with joined data
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.JoinRequest, error) {
	query := `
		SELECT jr.id, jr.full_name, jr.full_name_bn, jr.phone, jr.nid, jr.date_of_birth, 
		       jr.gender, jr.blood_group, jr.occupation, jr.address, jr.jurisdiction_id, 
		       jr.applied_at, jr.status, jr.referred_by_id, jr.rejection_reason, 
		       jr.processed_by_id, jr.created_at, jr.updated_at,
		       j.name as jurisdiction_name, u.full_name as referrer_name
		FROM join_requests jr
		JOIN jurisdictions j ON jr.jurisdiction_id = j.id
		LEFT JOIN users u ON jr.referred_by_id = u.id
		WHERE jr.id = $1
	`
	var jr models.JoinRequest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&jr.ID, &jr.FullName, &jr.FullNameBn, &jr.Phone, &jr.NID, &jr.DateOfBirth,
		&jr.Gender, &jr.BloodGroup, &jr.Occupation, &jr.Address, &jr.JurisdictionID,
		&jr.AppliedAt, &jr.Status, &jr.ReferredByID, &jr.RejectionReason,
		&jr.ProcessedByID, &jr.CreatedAt, &jr.UpdatedAt,
		&jr.JurisdictionName, &jr.ReferrerName,
	)
	if err != nil {
		return nil, err
	}
	return &jr, nil
}

// List returns join requests for a jurisdiction
func (r *Repository) List(ctx context.Context, jurisdictionID uuid.UUID, status string, limit, offset int) ([]*models.JoinRequest, error) {
	query := `
		SELECT jr.id, jr.full_name, jr.full_name_bn, jr.phone, jr.status, jr.applied_at,
		       j.name as jurisdiction_name
		FROM join_requests jr
		JOIN jurisdictions j ON jr.jurisdiction_id = j.id
		WHERE (jr.jurisdiction_id = $1 OR j.path <@ (SELECT path FROM jurisdictions WHERE id = $1))
	`
	args := []interface{}{jurisdictionID}
	nextArg := 2

	if status != "" {
		query += fmt.Sprintf(" AND jr.status = $%d", nextArg)
		args = append(args, status)
		nextArg++
	}

	query += fmt.Sprintf(" ORDER BY jr.applied_at DESC LIMIT $%d OFFSET $%d", nextArg, nextArg+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.JoinRequest
	for rows.Next() {
		var jr models.JoinRequest
		if err := rows.Scan(&jr.ID, &jr.FullName, &jr.FullNameBn, &jr.Phone, &jr.Status, &jr.AppliedAt, &jr.JurisdictionName); err != nil {
			return nil, err
		}
		list = append(list, &jr)
	}
	return list, nil
}

// UpdateStatus changes the status of a request and logs it in a transaction
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus, reason string, actorID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Get old status
	var oldStatus string
	err = tx.QueryRow(ctx, "SELECT status FROM join_requests WHERE id = $1", id).Scan(&oldStatus)
	if err != nil {
		return err
	}

	// 2. Update request
	query := `
		UPDATE join_requests 
		SET status = $1, rejection_reason = $2, processed_by_id = $3, updated_at = NOW() 
		WHERE id = $4
	`
	_, err = tx.Exec(ctx, query, newStatus, reason, actorID, id)
	if err != nil {
		return err
	}

	// 3. Log change
	logQuery := `
		INSERT INTO join_request_logs (request_id, actor_id, action, old_status, new_status, note)
		VALUES ($1, $2, 'status_change', $3, $4, $5)
	`
	_, err = tx.Exec(ctx, logQuery, id, actorID, oldStatus, newStatus, reason)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
