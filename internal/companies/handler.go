package companies

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Handler handles HTTP requests for companies
type Handler struct {
	service *Service
	logger  zerolog.Logger
}

// NewHandler creates a new company handler
func NewHandler(service *Service, logger zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers company routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Put("/{id}", h.Update)
	r.Delete("/{id}", h.Delete)
}

// List handles GET /companies
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := h.service.List(r.Context())
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list companies")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

// GetByID handles GET /companies/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	company, err := h.service.GetByID(r.Context(), id)
	if err == ErrNotFound {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get company")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// Create handles POST /companies
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		req.Status = "active"
	}

	company, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create company")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(company)
}

// Update handles PUT /companies/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	var req UpdateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	company, err := h.service.Update(r.Context(), id, req)
	if err == ErrNotFound {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to update company")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// Delete handles DELETE /companies/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err == ErrNotFound {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	} else if err != nil {
		h.logger.Error().Err(err).Msg("Failed to delete company")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
