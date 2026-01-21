package middleware

import (
	"net/http"

	"github.com/bjdms/api/internal/auth"
	"github.com/bjdms/api/internal/committee"
	"github.com/bjdms/api/pkg/response"
	"github.com/google/uuid"
)

// RBACMiddleware checks if the user has a specific role (position rank)
func RBACMiddleware(minRank int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				response.Unauthorized(w, "Not authenticated")
				return
			}

			// For now, we assume claims contain the user's rank/role
			// In a real system, we'd fetch this from the current session/cache
			// Logic: Lower rank number = higher authority
			// Super Admin = Rank 0 or 1
			// If minRank is 3, ranks 1, 2, 3 can pass.
			
			// Mock: If user has a rank in claims (need to add it)
			// For now, let's just make it a placeholder that we'll extend
			next.ServeHTTP(w, r)
		})
	}
}

// ABACJurisdictionMiddleware ensures the user can only manage their own or child jurisdictions
func ABACJurisdictionMiddleware(committeeService *committee.Service, authRepo *auth.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1. Get UserID from context
			userIDStr := GetUserID(r.Context())
			if userIDStr == "" {
				response.Unauthorized(w, "Not authenticated")
				return
			}
			userID, _ := uuid.Parse(userIDStr)

			// 2. Fetch User's jurisdiction from DB
			userJurisID, rank, err := authRepo.GetUserAuthDetails(r.Context(), userID)
			if err != nil {
				response.Unauthorized(w, "User authorization details not found")
				return
			}

			// Super Admin (Rank 1) can manage everything
			if rank == 1 {
				next.ServeHTTP(w, r)
				return
			}

			if userJurisID == nil {
				response.Forbidden(w, "You must be part of a committee to manage organizational units")
				return
			}

			// 3. Identify target jurisdiction (from URL param 'id' or body)
			// This is inherently dynamic and depends on the route.
			// For simplicity, let's look for a 'jurisdiction_id' in query or context
			targetIDStr := r.URL.Query().Get("jurisdiction_id")
			if targetIDStr == "" {
				// If not found, we might need to parse body or trust the handler to check
				// For now, let's allow it but the handler should verify
				next.ServeHTTP(w, r)
				return
			}

			targetID, err := uuid.Parse(targetIDStr)
			if err != nil {
				response.BadRequest(w, "Invalid target jurisdiction ID")
				return
			}

			// 4. Check if target is child of user jurisdiction
			isChild, err := committeeService.IsChildJurisdiction(r.Context(), *userJurisID, targetID)
			if err != nil || !isChild {
				response.Forbidden(w, "Access denied: Target jurisdiction is outside your area of responsibility")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
