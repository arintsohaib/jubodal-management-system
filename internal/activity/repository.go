package activity

import (
	"context"
	"fmt"
	"time"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for activities and tasks
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new activity repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// ACTIVITIES

// CreateActivity inserts a new activity record
func (r *Repository) CreateActivity(ctx context.Context, a *models.Activity) error {
	query := `
		INSERT INTO activities (user_id, jurisdiction_id, committee_id, title, description, category, activity_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	if a.ActivityDate.IsZero() {
		a.ActivityDate = time.Now()
	}
	return r.db.QueryRow(ctx, query,
		a.UserID, a.JurisdictionID, a.CommitteeID, a.Title, a.Description, a.Category, a.ActivityDate,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

// GetActivity retrieves an activity by ID
func (r *Repository) GetActivity(ctx context.Context, id uuid.UUID) (*models.Activity, error) {
	query := `
		SELECT a.id, a.user_id, a.jurisdiction_id, a.committee_id, a.title, a.description, a.category, a.activity_date, a.created_at, a.updated_at,
		       u.full_name as user_name, j.name as jurisdiction_name
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN jurisdictions j ON a.jurisdiction_id = j.id
		WHERE a.id = $1 AND a.deleted_at IS NULL
	`
	var a models.Activity
	err := r.db.QueryRow(ctx, query, id).Scan(
		&a.ID, &a.UserID, &a.JurisdictionID, &a.CommitteeID, &a.Title, &a.Description, &a.Category, &a.ActivityDate, &a.CreatedAt, &a.UpdatedAt,
		&a.UserName, &a.JurisdictionName,
	)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("activity not found")
	}
	return &a, err
}

// ListActivities returns activities filtered by jurisdiction and/or user
func (r *Repository) ListActivities(ctx context.Context, jurisdictionID *uuid.UUID, userID *uuid.UUID, limit, offset int) ([]*models.Activity, error) {
	query := `
		SELECT a.id, a.user_id, a.jurisdiction_id, a.committee_id, a.title, a.description, a.category, a.activity_date, a.created_at, a.updated_at,
		       u.full_name as user_name, j.name as jurisdiction_name
		FROM activities a
		JOIN users u ON a.user_id = u.id
		JOIN jurisdictions j ON a.jurisdiction_id = j.id
		WHERE a.deleted_at IS NULL
	`
	args := []interface{}{}
	if jurisdictionID != nil {
		args = append(args, *jurisdictionID)
		query += fmt.Sprintf(" AND a.jurisdiction_id = $%d", len(args))
	}
	if userID != nil {
		args = append(args, *userID)
		query += fmt.Sprintf(" AND a.user_id = $%d", len(args))
	}

	query += fmt.Sprintf(" ORDER BY a.activity_date DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Activity
	for rows.Next() {
		var a models.Activity
		err := rows.Scan(
			&a.ID, &a.UserID, &a.JurisdictionID, &a.CommitteeID, &a.Title, &a.Description, &a.Category, &a.ActivityDate, &a.CreatedAt, &a.UpdatedAt,
			&a.UserName, &a.JurisdictionName,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &a)
	}
	return list, nil
}

// TASKS

// CreateTask inserts a new task
func (r *Repository) CreateTask(ctx context.Context, t *models.Task) error {
	query := `
		INSERT INTO tasks (creator_id, assignee_id, committee_id, jurisdiction_id, title, description, status, priority, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		t.CreatorID, t.AssigneeID, t.CommitteeID, t.JurisdictionID, t.Title, t.Description, t.Status, t.Priority, t.DueDate,
	).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
}

// UpdateTaskStatus updates a task status
func (r *Repository) UpdateTaskStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
		UPDATE tasks 
		SET status = $1, 
		    updated_at = NOW(),
		    completed_at = CASE WHEN $1 = 'completed' THEN NOW() ELSE completed_at END
		WHERE id = $2 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

// ListTasks returns tasks filtered by jurisdiction or assignee
func (r *Repository) ListTasks(ctx context.Context, jurisdictionID *uuid.UUID, assigneeID *uuid.UUID, committeeID *uuid.UUID) ([]*models.Task, error) {
	query := `
		SELECT id, creator_id, assignee_id, committee_id, jurisdiction_id, title, description, status, priority, due_date, completed_at, verified_at, created_at, updated_at
		FROM tasks
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	if jurisdictionID != nil {
		args = append(args, *jurisdictionID)
		query += fmt.Sprintf(" AND jurisdiction_id = $%d", len(args))
	}
	if assigneeID != nil {
		args = append(args, *assigneeID)
		query += fmt.Sprintf(" AND assignee_id = $%d", len(args))
	}
	if committeeID != nil {
		args = append(args, *committeeID)
		query += fmt.Sprintf(" AND committee_id = $%d", len(args))
	}

	query += " ORDER BY priority ASC, due_date DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID, &t.CreatorID, &t.AssigneeID, &t.CommitteeID, &t.JurisdictionID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.DueDate, &t.CompletedAt, &t.VerifiedAt, &t.CreatedAt, &t.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, nil
}

// EVENTS

// CreateEvent inserts a new event
func (r *Repository) CreateEvent(ctx context.Context, e *models.Event) error {
	query := `
		INSERT INTO events (jurisdiction_id, organizer_id, title, description, location, start_time, end_time, is_public)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		e.JurisdictionID, e.OrganizerID, e.Title, e.Description, e.Location, e.StartTime, e.EndTime, e.IsPublic,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt)
}

// ListEvents returns events for a jurisdiction
func (r *Repository) ListEvents(ctx context.Context, jurisdictionID *uuid.UUID) ([]*models.Event, error) {
	query := `
		SELECT id, jurisdiction_id, organizer_id, title, description, location, start_time, end_time, is_public, created_at, updated_at
		FROM events
		WHERE deleted_at IS NULL
	`
	args := []interface{}{}
	if jurisdictionID != nil {
		args = append(args, *jurisdictionID)
		query += " AND jurisdiction_id = $1"
	}
	query += " ORDER BY start_time DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(
			&e.ID, &e.JurisdictionID, &e.OrganizerID, &e.Title, &e.Description, &e.Location, &e.StartTime, &e.EndTime, &e.IsPublic, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	return list, nil
}

// MarkAttendance records a user's presence at an event
func (r *Repository) MarkAttendance(ctx context.Context, eventID, userID uuid.UUID) error {
	query := `
		INSERT INTO event_attendance (event_id, user_id, attended_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (event_id, user_id) DO NOTHING
	`
	_, err := r.db.Exec(ctx, query, eventID, userID)
	return err
}
