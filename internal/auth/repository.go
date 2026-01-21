package auth

import (
	"context"
	"errors"
	"time"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// Repository defines database operations for authentication
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new auth repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// GetByPhone retrieves a user by their phone number
func (r *Repository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	query := `
		SELECT id, full_name, full_name_bn, phone, email, password_hash, 
		       is_active, verified_at, failed_login_attempts, locked_until, 
		       created_at, updated_at
		FROM users
		WHERE phone = $1 AND deleted_at IS NULL
	`

	var user models.User
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&user.ID, &user.FullName, &user.FullNameBn, &user.Phone, &user.Email, &user.PasswordHash,
		&user.IsActive, &user.VerifiedAt, &user.FailedLoginAttempts, &user.LockedUntil,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// IncrementFailedAttempts increases failed login counter and locks account if threshold reached
func (r *Repository) IncrementFailedAttempts(ctx context.Context, phone string, maxAttempts int, lockoutDuration time.Duration) error {
	query := `
		UPDATE users
		SET failed_login_attempts = failed_login_attempts + 1,
		    locked_until = CASE 
		        WHEN failed_login_attempts + 1 >= $2 THEN NOW() + $3::interval
		        ELSE locked_until 
		    END,
		    updated_at = NOW()
		WHERE phone = $1 AND deleted_at IS NULL
	`

	res, err := r.db.Exec(ctx, query, phone, maxAttempts, lockoutDuration.String())
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ResetFailedAttempts clears the failed login counter
func (r *Repository) ResetFailedAttempts(ctx context.Context, phone string) error {
	query := `
		UPDATE users
		SET failed_login_attempts = 0,
		    locked_until = NULL,
		    updated_at = NOW()
		WHERE phone = $1 AND deleted_at IS NULL
	`

	res, err := r.db.Exec(ctx, query, phone)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}

// CreateAuditLog inserts a new system audit event
func (r *Repository) CreateAuditLog(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (
			user_id, action, entity, entity_id, old_value, new_value, 
			ip_address, user_agent, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
	`
	_, err := r.db.Exec(ctx, query,
		log.UserID, log.Action, log.Entity, log.EntityID, log.OldValue, log.NewValue,
		log.IPAddress, log.UserAgent, log.Metadata,
	)
	return err
}

// Create inserts a new user record
func (r *Repository) Create(ctx context.Context, u *models.User) error {
	query := `
		INSERT INTO users (full_name, phone, nid, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		u.FullName, u.Phone, u.NID, u.PasswordHash, u.IsActive,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

// GetUserAuthDetails retrieves jurisdiction and rank for a user to support ABAC/RBAC
func (r *Repository) GetUserAuthDetails(ctx context.Context, userID uuid.UUID) (jurisdictionID *uuid.UUID, rank int, err error) {
	query := `
		SELECT jurisdiction_id, COALESCE(current_position_id, 999)
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	err = r.db.QueryRow(ctx, query, userID).Scan(&jurisdictionID, &rank)
	if err == pgx.ErrNoRows {
		return nil, 999, ErrUserNotFound
	}
	return jurisdictionID, rank, err
}
