package activity

import (
	"context"
	"fmt"
	"time"

	"github.com/bjdms/api/internal/models"
	"github.com/bjdms/api/internal/notification"
	"github.com/google/uuid"
)

// Service defines business logic for activities and tasks
type Service struct {
	repo         *Repository
	notification *notification.Service
}

// NewService creates a new activity service
func NewService(repo *Repository, ns *notification.Service) *Service {
	return &Service{repo: repo, notification: ns}
}

// ACTIVITIES

// LogActivity handles the creation of a new activity
func (s *Service) LogActivity(ctx context.Context, a *models.Activity) error {
	// 1. Validation logic (can user log in this jurisdiction?)
	// 2. Set defaults
	if a.Category == "" {
		a.Category = models.CategoryOrganizational
	}
	if a.ActivityDate.IsZero() {
		a.ActivityDate = time.Now()
	}

	return s.repo.CreateActivity(ctx, a)
}

// ListActivities returns activities with jurisdiction-aware filtering
func (s *Service) ListActivities(ctx context.Context, jurisdictionID *uuid.UUID, userID *uuid.UUID, page, pageSize int) ([]*models.Activity, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	return s.repo.ListActivities(ctx, jurisdictionID, userID, pageSize, offset)
}

// TASKS

// CreateTask handles task creation and assignment logic
func (s *Service) CreateTask(ctx context.Context, t *models.Task) error {
	// 1. Validation logic (e.g. Creator must have authority over assignee/jurisdiction)
	if t.Status == "" {
		t.Status = models.TaskStatusPending
	}
	if t.Priority == 0 {
		t.Priority = 3 // Medium default
	}

	err := s.repo.CreateTask(ctx, t)
	if err != nil {
		return err
	}

	// 3. Notify Assignee
	if t.AssigneeID != nil {
		s.notification.Create(ctx, &notification.Notification{
			UserID:         *t.AssigneeID,
			Type:           notification.TypeTaskAssigned,
			Title:          "New Task Assigned",
			Message:        fmt.Sprintf("You have been assigned a new task: %s", t.Title),
			JurisdictionID: t.JurisdictionID,
		})
	}

	return nil
}

// ListTasks returns tasks filtered by jurisdiction, assignee, or committee
func (s *Service) ListTasks(ctx context.Context, jurisdictionID *uuid.UUID, assigneeID *uuid.UUID, committeeID *uuid.UUID) ([]*models.Task, error) {
	return s.repo.ListTasks(ctx, jurisdictionID, assigneeID, committeeID)
}

// UpdateTaskStatus changes status and handles completion timestamps
func (s *Service) UpdateTaskStatus(ctx context.Context, taskID uuid.UUID, status string) error {
	// Valid status check
	validStatuses := map[string]bool{
		models.TaskStatusPending:    true,
		models.TaskStatusInProgress: true,
		models.TaskStatusCompleted:  true,
		models.TaskStatusVerified:   true,
		models.TaskStatusCancelled:  true,
	}

	if !validStatuses[status] {
		return fmt.Errorf("invalid task status: %s", status)
	}

	return s.repo.UpdateTaskStatus(ctx, taskID, status)
}

// EVENTS

// CreateEvent handles event creation logic
func (s *Service) CreateEvent(ctx context.Context, e *models.Event) error {
	if e.StartTime.IsZero() {
		return fmt.Errorf("start time is required")
	}
	return s.repo.CreateEvent(ctx, e)
}

// ListEvents returns events for a jurisdiction
func (s *Service) ListEvents(ctx context.Context, jurisdictionID *uuid.UUID) ([]*models.Event, error) {
	return s.repo.ListEvents(ctx, jurisdictionID)
}

// MarkAttendance records attendance at an event
func (s *Service) MarkAttendance(ctx context.Context, eventID, userID uuid.UUID) error {
	return s.repo.MarkAttendance(ctx, eventID, userID)
}
