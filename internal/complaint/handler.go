package complaint

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

// Handler handles HTTP requests for complaints
type Handler struct {
	service *Service
}

// NewHandler creates a new complaint handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes defines routes for complaints
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// Public Routes
	r.Group(func(r chi.Router) {
		r.Post("/public", h.SubmitAnonymous)
		r.Get("/status/{tracking_id}", h.CheckStatus)
	})

	// Protected Management Routes
	r.Group(func(r chi.Router) {
		// Middleware is applied in main.go (Auth + Jurisdiction ABAC)
		r.Get("/", h.ListComplaints)
		r.Get("/{tracking_id}", h.GetDetailed)
		r.Patch("/{id}/status", h.UpdateStatus)
	})

	return r
}

// SubmitAnonymous handles POST /api/v1/complaints/public
func (h *Handler) SubmitAnonymous(w http.ResponseWriter, r *http.Request) {
	var c models.Complaint
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	// For public submissions, ensure it's anonymous unless we'll allow auth later
	c.IsAnonymous = true
	c.UserID = nil

	// Get IP for the service to hash
	ip := r.RemoteAddr

	if err := h.service.SubmitComplaint(r.Context(), &c, ip); err != nil {
		response.InternalError(w, "Failed to submit complaint", "")
		return
	}

	response.Created(w, map[string]string{
		"tracking_id": c.TrackingID,
		"status":      c.Status,
	}, "Complaint submitted successfully. Please save your tracking ID.")
}

// CheckStatus handles GET /api/v1/complaints/status/{tracking_id}
func (h *Handler) CheckStatus(w http.ResponseWriter, r *http.Request) {
	trackingID := chi.URLParam(r, "tracking_id")
	if trackingID == "" {
		response.BadRequest(w, "Tracking ID is required")
		return
	}

	c, err := h.service.GetComplaintStatus(r.Context(), trackingID)
	if err != nil {
		response.NotFound(w, "Complaint not found")
		return
	}

	// We only return public info
	pubInfo := map[string]interface{}{
		"tracking_id": c.TrackingID,
		"status":      c.Status,
		"created_at":  c.CreatedAt,
		"subject":     c.Subject,
	}
	if c.ResolutionNotes != nil {
		pubInfo["resolution_notes"] = *c.ResolutionNotes
	}

	response.Success(w, pubInfo, "")
}

// ListComplaints handles GET /api/v1/complaints
func (h *Handler) ListComplaints(w http.ResponseWriter, r *http.Request) {
	jurisIDStr := r.URL.Query().Get("jurisdiction_id")
	if jurisIDStr == "" {
		response.BadRequest(w, "jurisdiction_id is required")
		return
	}
	jurisID, err := uuid.Parse(jurisIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid jurisdiction_id")
		return
	}

	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	list, err := h.service.ListJurisdictionComplaints(r.Context(), jurisID, status, page, pageSize)
	if err != nil {
		response.InternalError(w, "Failed to fetch complaints", "")
		return
	}

	response.Success(w, list, "")
}

// GetDetailed handles GET /api/v1/complaints/{tracking_id}
func (h *Handler) GetDetailed(w http.ResponseWriter, r *http.Request) {
	trackingID := chi.URLParam(r, "tracking_id")
	c, err := h.service.repo.GetByTrackingID(r.Context(), trackingID)
	if err != nil {
		response.NotFound(w, "Complaint not found")
		return
	}

	response.Success(w, c, "")
}

// UpdateStatus handles PATCH /api/v1/complaints/{id}/status
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid complaint ID")
		return
	}

	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := uuid.Parse(userIDStr)

	var req struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.UpdateComplaintStatus(r.Context(), id, userID, req.Status, req.Note); err != nil {
		response.InternalError(w, "Failed to update status", "")
		return
	}

	response.Success(w, nil, "Status updated successfully")
}
