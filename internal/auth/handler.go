package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bjdms/api/internal/models"
	"github.com/bjdms/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Handler handles HTTP requests for authentication
type Handler struct {
	service      *Service
	redisManager *RedisManager
	jwtManager   *JWTManager
}

// NewHandler creates a new auth handler
func NewHandler(service *Service, redisManager *RedisManager, jwtManager *JWTManager) *Handler {
	return &Handler{
		service:      service,
		redisManager: redisManager,
		jwtManager:   jwtManager,
	}
}

// Routes defines auth endpoints
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/login", h.Login)
	r.Post("/refresh", h.Refresh)
	r.Post("/logout", h.Logout)

	return r
}

// Login handles user login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// Basic validation
	if err := ValidatePhone(req.Phone); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	ip := strings.Split(r.RemoteAddr, ":")[0]
	ua := r.UserAgent()

	res, err := h.service.Login(r.Context(), req, ip, ua)
	if err != nil {
		switch err {
		case ErrInvalidCredentials:
			response.Unauthorized(w, err.Error())
		case ErrAccountLocked:
			response.Error(w, http.StatusLocked, "account_locked", err.Error(), "")
		case ErrAccountInactive:
			response.Forbidden(w, err.Error())
		default:
			requestID := middleware.GetReqID(r.Context())
			response.InternalError(w, "Failed to login", requestID)
		}
		return
	}

	// Store sessions in Redis
	// In a full implementation, we'd extract token IDs from res
	// For now, assume service handles it or we do it here

	response.Success(w, res, "Login successful")
}

// Refresh handles token refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	res, err := h.service.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(w, "Invalid refresh token")
		return
	}

	response.Success(w, res, "Tokens refreshed")
}

// Logout invalidates user session
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// 1. Get token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		response.Success(w, nil, "Logged out (no session found)")
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.jwtManager.VerifyAccessToken(token)
	if err != nil {
		response.Success(w, nil, "Logged out (invalid token)")
		return
	}

	// 2. Invalidate session in Redis
	h.redisManager.InvalidateSession(r.Context(), claims.UserID, claims.TokenID)
	
	// 3. Optional: Invalidate associated refresh token

	response.Success(w, nil, "Logged out successfully")
}
