package activity

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

// Handler handles HTTP requests for activities and tasks
type Handler struct {
	service *Service
}

// NewHandler creates a new activity handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes defines routes for activities and tasks
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	// Activities
	r.Post("/", h.LogActivity)
	r.Get("/", h.ListActivities)
	r.Get("/{id}", h.GetActivity)

	// Tasks
	r.Post("/tasks", h.CreateTask)
	r.Get("/tasks", h.ListTasks)
	r.Patch("/tasks/{id}/status", h.UpdateTaskStatus)

	// Events
	r.Post("/events", h.CreateEvent)
	r.Get("/events", h.ListEvents)
	r.Post("/events/{id}/attendance", h.MarkAttendance)

	return r
}

// ... existing methods ...

// EVENTS

// CreateEvent handles POST /api/v1/activities/events
func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	creatorIDStr := middleware.GetUserID(r.Context())
	creatorID, _ := uuid.Parse(creatorIDStr)

	var e models.Event
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	e.OrganizerID = creatorID

	if err := h.service.CreateEvent(r.Context(), &e); err != nil {
		response.InternalError(w, "Failed to create event", "")
		return
	}

	response.Created(w, e, "Event created successfully")
}

// ListEvents handles GET /api/v1/activities/events
func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	jurisIDStr := r.URL.Query().Get("jurisdiction_id")
	var jurisID *uuid.UUID
	if id, err := uuid.Parse(jurisIDStr); err == nil {
		jurisID = &id
	}

	list, err := h.service.ListEvents(r.Context(), jurisID)
	if err != nil {
		response.InternalError(w, "Failed to fetch events", "")
		return
	}

	response.Success(w, list, "")
}

// MarkAttendance handles POST /api/v1/activities/events/{id}/attendance
func (h *Handler) MarkAttendance(w http.ResponseWriter, r *http.Request) {
	eventIDStr := chi.URLParam(r, "id")
	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		response.BadRequest(w, "Invalid event ID")
		return
	}

	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := uuid.Parse(userIDStr)

	if err := h.service.MarkAttendance(r.Context(), eventID, userID); err != nil {
		response.InternalError(w, "Failed to mark attendance", "")
		return
	}

	response.Created(w, nil, "Attendance marked")
}

// ACTIVITIES

// LogActivity handles POST /api/v1/activities
func (h *Handler) LogActivity(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := uuid.Parse(userIDStr)

	var a models.Activity
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	a.UserID = userID

	if err := h.service.LogActivity(r.Context(), &a); err != nil {
		response.InternalError(w, "Failed to log activity", "")
		return
	}

	response.Created(w, a, "Activity logged successfully")
}

// ListActivities handles GET /api/v1/activities
func (h *Handler) ListActivities(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	
	jurisIDStr := r.URL.Query().Get("jurisdiction_id")
	var jurisID *uuid.UUID
	if id, err := uuid.Parse(jurisIDStr); err == nil {
		jurisID = &id
	}

	list, err := h.service.ListActivities(r.Context(), jurisID, nil, page, pageSize)
	if err != nil {
		response.InternalError(w, "Failed to fetch activities", "")
		return
	}

	response.Success(w, list, "")
}

// GetActivity handles GET /api/v1/activities/{id}
func (h *Handler) GetActivity(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid activity ID")
		return
	}

	a, err := h.service.repo.GetActivity(r.Context(), id)
	if err != nil {
		response.NotFound(w, "Activity not found")
		return
	}

	response.Success(w, a, "")
}

// TASKS

// CreateTask handles POST /api/v1/activities/tasks
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	creatorIDStr := middleware.GetUserID(r.Context())
	creatorID, _ := uuid.Parse(creatorIDStr)

	var t models.Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	t.CreatorID = creatorID

	if err := h.service.CreateTask(r.Context(), &t); err != nil {
		response.InternalError(w, "Failed to create task", "")
		return
	}

	response.Created(w, t, "Task created successfully")
}

// ListTasks handles GET /api/v1/activities/tasks
func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	jurisIDStr := r.URL.Query().Get("jurisdiction_id")
	var jurisID *uuid.UUID
	if id, err := uuid.Parse(jurisIDStr); err == nil {
		jurisID = &id
	}

	list, err := h.service.ListTasks(r.Context(), jurisID, nil, nil)
	if err != nil {
		response.InternalError(w, "Failed to fetch tasks", "")
		return
	}

	response.Success(w, list, "")
}

// UpdateTaskStatus handles PATCH /api/v1/activities/tasks/{id}/status
func (h *Handler) UpdateTaskStatus(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(w, "Invalid task ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}

	if err := h.service.UpdateTaskStatus(r.Context(), id, req.Status); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Success(w, nil, "Task status updated")
}
