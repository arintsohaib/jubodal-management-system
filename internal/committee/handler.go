package committee

import (
	"encoding/json"
	"net/http"

	"github.com/bjdms/api/internal/models"
	"github.com/bjdms/api/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler handles HTTP requests for committees and jurisdictions
type Handler struct {
	service *Service
}

// NewHandler creates a new committee handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes defines routes for committees and jurisdictions
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// Jurisdictions
	r.Post("/jurisdictions", h.CreateJurisdiction)
	r.Get("/jurisdictions", h.ListJurisdictions)

	// Committees
	r.Post("/committees", h.CreateCommittee)
	r.Get("/committees/{id}/members", h.ListMembers)
	r.Post("/committees/{id}/members", h.AddMember)

	return r
}

// CreateJurisdiction handles POST /jurisdictions
func (h *Handler) CreateJurisdiction(w http.ResponseWriter, r *http.Request) {
	var j models.Jurisdiction
	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.CreateJurisdiction(r.Context(), &j); err != nil {
		response.Error(w, http.StatusInternalServerError, "jurisdiction_error", err.Error(), "")
		return
	}

	response.Success(w, j, "Jurisdiction created successfully")
}

// ListJurisdictions handles GET /jurisdictions
func (h *Handler) ListJurisdictions(w http.ResponseWriter, r *http.Request) {
	parentStr := r.URL.Query().Get("parent_id")
	var parentID *uuid.UUID
	if parentStr != "" {
		id, err := uuid.Parse(parentStr)
		if err == nil {
			parentID = &id
		}
	}

	list, err := h.service.ListJurisdictionTree(r.Context(), parentID)
	if err != nil {
		response.InternalError(w, "Failed to fetch jurisdictions", "")
		return
	}

	response.Success(w, list, "")
}

// CreateCommittee handles POST /committees
func (h *Handler) CreateCommittee(w http.ResponseWriter, r *http.Request) {
	var c models.Committee
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.CreateCommittee(r.Context(), &c); err != nil {
		response.Error(w, http.StatusConflict, "committee_error", err.Error(), "")
		return
	}

	response.Success(w, c, "Committee created successfully")
}

// AddMember handles POST /committees/{id}/members
func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request) {
	committeeIDStr := chi.URLParam(r, "id")
	committeeID, err := uuid.Parse(committeeIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid committee ID")
		return
	}

	var m models.CommitteeMember
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	m.CommitteeID = committeeID

	if err := h.service.AddMember(r.Context(), &m); err != nil {
		response.Error(w, http.StatusBadRequest, "member_error", err.Error(), "")
		return
	}

	response.Success(w, m, "Member added successfully")
}

// ListMembers handles GET /committees/{id}/members
func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	committeeIDStr := chi.URLParam(r, "id")
	committeeID, err := uuid.Parse(committeeIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid committee ID")
		return
	}

	members, err := h.service.repo.GetCommitteeMembers(r.Context(), committeeID)
	if err != nil {
		response.InternalError(w, "Failed to fetch members", "")
		return
	}

	response.Success(w, members, "")
}
