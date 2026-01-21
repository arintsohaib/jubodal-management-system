package join

import (
	"context"
	"fmt"
	"time"

	"github.com/bjdms/api/internal/auth"
	"github.com/bjdms/api/internal/models"
	"github.com/bjdms/api/internal/notification"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service handles business logic for joining Jubodal
type Service struct {
	repo         *Repository
	userRepo     *auth.Repository
	notification *notification.Service
}

// NewService creates a new join service
func NewService(repo *Repository, userRepo *auth.Repository, ns *notification.Service) *Service {
	return &Service{repo: repo, userRepo: userRepo, notification: ns}
}

// SubmitApplication handles public join request submission
func (s *Service) SubmitApplication(ctx context.Context, jr *models.JoinRequest) error {
	// 1. Age Validation (18-40 is typically Jubo Dal bracket, allowing some flexibility)
	age := time.Since(jr.DateOfBirth).Hours() / 24 / 365
	if age < 18 {
		return fmt.Errorf("applicant must be at least 18 years old")
	}

	// 2. Duplicate check (Phone/NID)
	// This would typically involve checking existing users and pending join requests
	// Simplified for this implementation

	jr.Status = models.JoinRequestStatusPending
	err := s.repo.Create(ctx, jr)
	if err != nil {
		return err
	}

	// Notify Jurisdiction Leaders of new application
	// s.notification.Create(ctx, &notification.Notification{...})

	return nil
}

// ApproveRequest approves a join request and creates a user account
func (s *Service) ApproveRequest(ctx context.Context, id uuid.UUID, actorID uuid.UUID) error {
	jr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if jr.Status != models.JoinRequestStatusPending && jr.Status != models.JoinRequestStatusUnderReview {
		return fmt.Errorf("request is in %s status and cannot be approved", jr.Status)
	}

	// 1. Create User Account
	tempPassword := "Jubodal@123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)

	user := &models.User{
		FullName:     jr.FullName,
		Phone:        jr.Phone,
		NID:          jr.NID,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user account: %w", err)
	}

	// 2. Update Join Request Status
	err = s.repo.UpdateStatus(ctx, id, models.JoinRequestStatusApproved, "Approved by committee", actorID)
	if err != nil {
		return err
	}

	// 3. Notify the successful applicant
	s.notification.Create(ctx, &notification.Notification{
		UserID:         user.ID,
		Type:           notification.TypeJoinRequest,
		Title:          "Welcome to Jubodal!",
		Message:        "Your membership application has been approved. You can now log in using your phone number.",
		JurisdictionID: jr.JurisdictionID,
	})

	return nil
}

// RejectRequest rejects an application with a reason
func (s *Service) RejectRequest(ctx context.Context, id uuid.UUID, reason string, actorID uuid.UUID) error {
	return s.repo.UpdateStatus(ctx, id, models.JoinRequestStatusRejected, reason, actorID)
}

// ListRequests returns applications for a leader's jurisdiction
func (s *Service) ListRequests(ctx context.Context, jurisdictionID uuid.UUID, status string, page, pageSize int) ([]*models.JoinRequest, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	return s.repo.List(ctx, jurisdictionID, status, pageSize, offset)
}
