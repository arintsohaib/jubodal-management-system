package join

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/bjdms/api/internal/middleware"
	"github.com/bjdms/api/internal/models"
	"github.com/bjdms/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles join request HTTP endpoints
type Handler struct {
	service *Service
}

// NewHandler creates a new join handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes returns routes for protected management
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
	r.Patch("/{id}/approve", h.Approve)
	r.Patch("/{id}/reject", h.Reject)

	return r
}

// PublicRoutes returns routes for anonymous application
func (h *Handler) PublicRoutes() chi.Router {
	r := chi.NewRouter()
	r.Post("/apply", h.Submit)
	return r
}

// Submit handles POST /api/v1/public/join/apply
func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	var jr models.JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&jr); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.SubmitApplication(r.Context(), &jr); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, jr, "Application submitted successfully. We will contact you after review.")
}

// List handles GET /api/v1/join-requests
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// ABAC middleware will have populated user details
	jurisIDStr := r.URL.Query().Get("jurisdiction_id") // Fallback or override
	if jurisIDStr == "" {
		// In a real system, we'd get the leader's jurisdiction from context
		response.BadRequest(w, "jurisdiction_id is required")
		return
	}
	jurisID, _ := uuid.Parse(jurisIDStr)

	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	list, err := h.service.ListRequests(r.Context(), jurisID, status, page, pageSize)
	if err != nil {
		response.InternalError(w, "Failed to list requests", "")
		return
	}

	response.Success(w, list, "")
}

// Get handles GET /api/v1/join-requests/{id}
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	jr, err := h.service.repo.GetByID(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Join request not found")
		return
	}
	response.Success(w, jr, "")
}

// Approve handles PATCH /api/v1/join-requests/{id}/approve
func (h *Handler) Approve(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	actorIDStr := middleware.GetUserID(r.Context())
	actorID, _ := uuid.Parse(actorIDStr)

	if err := h.service.ApproveRequest(r.Context(), id, actorID); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, nil, "Request approved and member account created")
}

// Reject handles PATCH /api/v1/join-requests/{id}/reject
func (h *Handler) Reject(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.Parse(chi.URLParam(r, "id"))
	actorIDStr := middleware.GetUserID(r.Context())
	actorID, _ := uuid.Parse(actorIDStr)

	var body struct {
		Reason string `json:"reason"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	if err := h.service.RejectRequest(r.Context(), id, body.Reason, actorID); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, nil, "Request rejected")
}
