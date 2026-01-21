package complaint

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/bjdms/api/internal/models"
	"github.com/bjdms/api/internal/notification"
	"github.com/google/uuid"
)

// Service handles business logic for complaints
type Service struct {
	repo         *Repository
	notification *notification.Service
}

// NewService creates a new complaint service
func NewService(repo *Repository, ns *notification.Service) *Service {
	return &Service{repo: repo, notification: ns}
}

// SubmitComplaint handles raw submission including tracking ID and IP hashing
func (s *Service) SubmitComplaint(ctx context.Context, c *models.Complaint, ipAddress string) error {
	// 1. Generate human-readable Tracking ID (C-YYYY-[RANDOM])
	c.TrackingID = generateTrackingID()

	// 2. Hash IP for anonymous tracking/spam control
	if c.IsAnonymous {
		hash := sha256.Sum256([]byte(ipAddress))
		hashStr := fmt.Sprintf("%x", hash)
		c.AnonymousIPHash = &hashStr
		c.ComplainantName = nil
		c.ComplainantContact = nil
	}

	// 3. Defaults
	if c.Status == "" {
		c.Status = models.ComplaintStatusReceived
	}

	err := s.repo.CreateComplaint(ctx, c)
	if err != nil {
		return err
	}

	// 4. Notify Jurisdiction Leaders
	// Realistically, we'd fetch leader IDs, but for this demo, we'll notify the 'Jurisdiction Feed'
	// System assumes an automated listener for these alerts.
	s.notification.Create(ctx, &notification.Notification{
		Type:           notification.TypeComplaintAlert,
		Title:          "New Complaint Filed",
		Message:        fmt.Sprintf("Tracking ID: %s. A new complaint has been submitted to your jurisdiction.", c.TrackingID),
		JurisdictionID: c.JurisdictionID,
	})

	return nil
}

// GetComplaintStatus allows public/anonymous lookup by tracking ID
func (s *Service) GetComplaintStatus(ctx context.Context, trackingID string) (*models.Complaint, error) {
	// We mask sensitive data for public lookup
	c, err := s.repo.GetByTrackingID(ctx, trackingID)
	if err != nil {
		return nil, err
	}

	if c.IsAnonymous {
		c.ComplainantName = nil
		c.ComplainantContact = nil
	}
	
	return c, nil
}

// UpdateComplaintStatus handles status transitions with audit trail
func (s *Service) UpdateComplaintStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status, note string) error {
	// Logic to validate status transition could be added here
	return s.repo.UpdateStatus(ctx, id, userID, status, note)
}

// ListJurisdictionComplaints returns complaints for authorized leaders
func (s *Service) ListJurisdictionComplaints(ctx context.Context, jurisdictionID uuid.UUID, status string, page, pageSize int) ([]*models.Complaint, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	return s.repo.ListComplaints(ctx, jurisdictionID, status, pageSize, offset)
}

// Helper: Generate a unique tracking ID
func generateTrackingID() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // Avoid ambiguous chars 0,O,1,I
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return fmt.Sprintf("C-%d-%s", now.Year(), string(code))
}
