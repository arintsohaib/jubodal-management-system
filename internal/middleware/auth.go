package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/bjdms/api/internal/auth"
	"github.com/bjdms/api/pkg/response"
)

type contextKey string

const (
	UserIDKey  contextKey = "user_id"
	ClaimsKey  contextKey = "claims"
)

// AuthMiddleware authenticates requests using JWT access tokens
func AuthMiddleware(jwtManager *auth.JWTManager, redisManager *auth.RedisManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "Missing authorization header")
				return
			}

			// 2. Parse Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				response.Unauthorized(w, "Invalid authorization header format")
				return
			}

			token := parts[1]

			// 3. Verify access token
			claims, err := jwtManager.VerifyAccessToken(token)
			if err != nil {
				response.Unauthorized(w, "Invalid or expired token")
				return
			}

			// 4. Verify session in Redis
			active, err := redisManager.ValidateSession(r.Context(), claims.UserID, claims.TokenID)
			if err != nil {
				// Log error but proceed if Redis is down? (Fail closed suggested for security)
				response.InternalError(w, "Failed to validate session", "")
				return
			}

			if !active {
				response.Unauthorized(w, "Session expired or revoked")
				return
			}

			// 5. Check if token is revoked
			revoked, err := redisManager.IsRevoked(r.Context(), claims.TokenID)
			if err == nil && revoked {
				response.Unauthorized(w, "Token revoked")
				return
			}

			// 6. Add claims and user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, ClaimsKey, claims)
			
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID retrieves the authenticated user ID from context
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value(UserIDKey).(string); ok {
		return id
	}
	return ""
}

// GetClaims retrieves the JWT claims from context
func GetClaims(ctx context.Context) *auth.Claims {
	if claims, ok := ctx.Value(ClaimsKey).(*auth.Claims); ok {
		return claims
	}
	return nil
}
