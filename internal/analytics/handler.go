package analytics

import (
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	client *Client
}

func NewHandler(client *Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/pulse", h.GetPulse)
	r.Post("/query", h.QueryAssistant)
	return r
}

func (h *Handler) QueryAssistant(w http.ResponseWriter, r *http.Request) {
	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := h.client.QueryAssistant(payload)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) GetPulse(w http.ResponseWriter, r *http.Request) {
	data, err := h.client.GetPulse()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Analytics service is currently unavailable",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
