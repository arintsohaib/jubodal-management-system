package auth

import (
	"context"
	"errors"

	"github.com/bjdms/api/config"
	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
)

var (
	ErrAccountLocked    = errors.New("account is temporarily locked due to multiple failed login attempts")
	ErrAccountInactive  = errors.New("account is inactive")
	ErrInvalidCredentials = errors.New("invalid phone number or password")
)

// Service defines business logic for authentication
type Service struct {
	repo         *Repository
	redisManager *RedisManager
	jwtManager   *JWTManager
	config       *config.Config
}

// NewService creates a new auth service
func NewService(repo *Repository, redisManager *RedisManager, jwtManager *JWTManager, cfg *config.Config) *Service {
	return &Service{
		repo:         repo,
		redisManager: redisManager,
		jwtManager:   jwtManager,
		config:       cfg,
	}
}

// Login authenticates a user and returns JWT tokens
func (s *Service) Login(ctx context.Context, req models.LoginRequest, ip, ua string) (*models.LoginResponse, error) {
	// 1. Get user by phone
	user, err := s.repo.GetByPhone(ctx, req.Phone)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. Check if account is active
	if !user.IsActive {
		s.logAuthEvent(ctx, user.ID, "login_failed", "inactive_account", ip, ua)
		return nil, ErrAccountInactive
	}

	// 3. Check if account is locked
	if user.IsLocked() {
		s.logAuthEvent(ctx, user.ID, "login_failed", "account_locked", ip, ua)
		return nil, ErrAccountLocked
	}

	// 4. Verify password
	if !ComparePassword(req.Password, user.PasswordHash) {
		// Increment failed attempts
		s.repo.IncrementFailedAttempts(ctx, req.Phone, s.config.MaxFailedAttempts, s.config.LockoutDuration)
		s.logAuthEvent(ctx, user.ID, "login_failed", "invalid_password", ip, ua)
		return nil, ErrInvalidCredentials
	}

	// 5. Successful login - reset failed attempts
	s.repo.ResetFailedAttempts(ctx, req.Phone)

	// 6. Generate tokens
	accessToken, accessTokenID, err := s.jwtManager.GenerateAccessToken(user.ID, user.Phone, user.IsVerified())
	if err != nil {
		return nil, err
	}

	refreshToken, refreshTokenID, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 7. Store sessions in Redis
	if err := s.redisManager.SetSession(ctx, user.ID.String(), accessTokenID, s.config.JWTAccessExpiry); err != nil {
		return nil, err
	}
	if err := s.redisManager.SetSession(ctx, user.ID.String(), refreshTokenID, s.config.JWTRefreshExpiry); err != nil {
		return nil, err
	}

	// 8. Audit log success
	s.logAuthEvent(ctx, user.ID, "login_success", "", ip, ua)

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.JWTAccessExpiry.Seconds()),
		User:         user,
	}, nil
}

// RefreshTokens validates refresh token and returns new tokens
func (s *Service) RefreshTokens(ctx context.Context, refreshToken string) (*models.LoginResponse, error) {
	// 1. Verify refresh token
	claims, err := s.jwtManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 2. Mock session/user check using claims
	uid, _ := uuid.Parse(claims.Subject)

	accessToken, _, _ := s.jwtManager.GenerateAccessToken(uid, "user", true)
	
	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.JWTAccessExpiry.Seconds()),
	}, nil
}

// logAuthEvent creates an audit log for authentication actions
func (s *Service) logAuthEvent(ctx context.Context, userID interface{}, action, reason string, ip, ua string) {
	metadata := make(map[string]interface{})
	if reason != "" {
		metadata["reason"] = reason
	}

	log := &models.AuditLog{
		Action:    action,
		Entity:    "user",
		IPAddress: &ip,
		UserAgent: &ua,
		Metadata:  metadata,
	}

	if id, ok := userID.(uuid.UUID); ok {
		log.UserID = &id
	}

	// Async log to not block request
	go func() {
		s.repo.CreateAuditLog(context.Background(), log)
	}()
}
