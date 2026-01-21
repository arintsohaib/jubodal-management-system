package complaint

import (
	"context"
	"fmt"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for complaints
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new complaint repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateComplaint inserts a new complaint
func (r *Repository) CreateComplaint(ctx context.Context, c *models.Complaint) error {
	query := `
		INSERT INTO complaints (
			tracking_id, user_id, jurisdiction_id, is_anonymous, 
			complainant_name, complainant_contact, subject, description, 
			status, anonymous_ip_hash
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		c.TrackingID, c.UserID, c.JurisdictionID, c.IsAnonymous,
		c.ComplainantName, c.ComplainantContact, c.Subject, c.Description,
		c.Status, c.AnonymousIPHash,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

// GetByTrackingID retrieves a complaint by its human-readable tracking ID
func (r *Repository) GetByTrackingID(ctx context.Context, trackingID string) (*models.Complaint, error) {
	query := `
		SELECT c.id, c.tracking_id, c.user_id, c.jurisdiction_id, c.is_anonymous, 
		       c.complainant_name, c.complainant_contact, c.subject, c.description, 
		       c.status, c.assigned_to_id, c.resolution_notes, c.created_at, c.updated_at,
		       j.name as jurisdiction_name
		FROM complaints c
		JOIN jurisdictions j ON c.jurisdiction_id = j.id
		WHERE c.tracking_id = $1 AND c.deleted_at IS NULL
	`
	var c models.Complaint
	err := r.db.QueryRow(ctx, query, trackingID).Scan(
		&c.ID, &c.TrackingID, &c.UserID, &c.JurisdictionID, &c.IsAnonymous,
		&c.ComplainantName, &c.ComplainantContact, &c.Subject, &c.Description,
		&c.Status, &c.AssignedToID, &c.ResolutionNotes, &c.CreatedAt, &c.UpdatedAt,
		&c.JurisdictionName,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("complaint not found")
	}
	return &c, err
}

// ListComplaints returns complaints filtered by jurisdiction and status
func (r *Repository) ListComplaints(ctx context.Context, jurisdictionID uuid.UUID, status string, limit, offset int) ([]*models.Complaint, error) {
	query := `
		SELECT c.id, c.tracking_id, c.user_id, c.jurisdiction_id, c.is_anonymous, 
		       c.complainant_name, c.complainant_contact, c.subject, c.description, 
		       c.status, c.assigned_to_id, c.resolution_notes, c.created_at, c.updated_at,
		       j.name as jurisdiction_name, u.full_name as assigned_to_name
		FROM complaints c
		JOIN jurisdictions j ON c.jurisdiction_id = j.id
		LEFT JOIN users u ON c.assigned_to_id = u.id
		WHERE c.jurisdiction_id = $1 AND c.deleted_at IS NULL
	`
	args := []interface{}{jurisdictionID}
	
	if status != "" {
		args = append(args, status)
		query += fmt.Sprintf(" AND c.status = $%d", len(args))
	}

	query += fmt.Sprintf(" ORDER BY c.created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Complaint
	for rows.Next() {
		var c models.Complaint
		err := rows.Scan(
			&c.ID, &c.TrackingID, &c.UserID, &c.JurisdictionID, &c.IsAnonymous,
			&c.ComplainantName, &c.ComplainantContact, &c.Subject, &c.Description,
			&c.Status, &c.AssignedToID, &c.ResolutionNotes, &c.CreatedAt, &c.UpdatedAt,
			&c.JurisdictionName, &c.AssignedToName,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, nil
}

// UpdateStatus changes the status and logs the action in a transaction
func (r *Repository) UpdateStatus(ctx context.Context, complaintID uuid.UUID, userID uuid.UUID, newStatus, note string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Get current status
	var oldStatus string
	err = tx.QueryRow(ctx, "SELECT status FROM complaints WHERE id = $1 FOR UPDATE", complaintID).Scan(&oldStatus)
	if err != nil {
		return err
	}

	// 2. Update status
	_, err = tx.Exec(ctx, "UPDATE complaints SET status = $1, updated_at = NOW() WHERE id = $2", newStatus, complaintID)
	if err != nil {
		return err
	}

	// 3. Log the change
	logQuery := `
		INSERT INTO complaint_logs (complaint_id, user_id, action, old_status, new_status, note)
		VALUES ($1, $2, 'status_change', $3, $4, $5)
	`
	_, err = tx.Exec(ctx, logQuery, complaintID, userID, oldStatus, newStatus, note)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// CreateEvidence links a file to a complaint
func (r *Repository) CreateEvidence(ctx context.Context, e *models.ComplaintEvidence) error {
	query := `INSERT INTO complaint_evidence (complaint_id, file_path, file_type) VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, e.ComplaintID, e.FilePath, e.FileType).Scan(&e.ID, &e.CreatedAt)
}
