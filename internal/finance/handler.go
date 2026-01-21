package finance

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

// Handler handles financial HTTP requests
type Handler struct {
	service *Service
}

// NewHandler creates a new finance handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Routes defines routes for financial management
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/categories", h.ListCategories)
	r.Post("/transactions", h.RecordTransaction)
	r.Get("/statement", h.GetStatement)

	return r
}

// ListCategories handles GET /api/v1/finance/categories
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	transType := r.URL.Query().Get("type")
	list, err := h.service.ListCategories(r.Context(), transType)
	if err != nil {
		response.InternalError(w, "Failed to fetch categories", "")
		return
	}
	response.Success(w, list, "")
}

// RecordTransaction handles POST /api/v1/finance/transactions
func (h *Handler) RecordTransaction(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := uuid.Parse(userIDStr)

	var t models.FinanceTransaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		response.BadRequest(w, "Invalid request body")
		return
	}
	t.UserID = userID

	if err := h.service.RecordTransaction(r.Context(), &t); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	response.Created(w, t, "Transaction recorded successfully")
}

// GetStatement handles GET /api/v1/finance/statement?jurisdiction_id=...
func (h *Handler) GetStatement(w http.ResponseWriter, r *http.Request) {
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

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	balance, transactions, err := h.service.GetJurisdictionStatement(r.Context(), jurisID, page, pageSize)
	if err != nil {
		response.InternalError(w, "Failed to fetch statement", "")
		return
	}

	response.Success(w, map[string]interface{}{
		"balance":      balance,
		"transactions": transactions,
	}, "")
}
