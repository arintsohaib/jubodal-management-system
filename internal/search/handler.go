package search

import (
	"net/http"
	"strconv"

	"github.com/bjdms/api/pkg/response"
	"github.com/go-chi/chi/v5"
)

// Handler handles search HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new search handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes defines routes for search
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.Search)
	r.Get("/autocomplete", h.Autocomplete)

	return r
}

// Search handles GET /api/v1/search?q=...
func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "Query parameter 'q' is required")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = 20
	}
	offset := page * pageSize

	results, err := h.service.SearchGlobal(r.Context(), query, pageSize, offset)
	if err != nil {
		response.InternalError(w, "Search failed", err.Error())
		return
	}

	response.Success(w, results, "")
}

// Autocomplete handles GET /api/v1/search/autocomplete?q=...
func (h *Handler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		response.BadRequest(w, "Query parameter 'q' is required")
		return
	}

	results, err := h.service.Autocomplete(r.Context(), query)
	if err != nil {
		response.InternalError(w, "Autocomplete failed", err.Error())
		return
	}

	response.Success(w, results, "")
}
