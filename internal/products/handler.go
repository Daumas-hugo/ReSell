package products

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Handler handles HTTP requests for products
type Handler struct {
	service *Service
	logger  zerolog.Logger
}

// NewHandler creates a new product handler
func NewHandler(service *Service, logger zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers product routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Get("/{id}/variants", h.ListVariants)
	r.Post("/{id}/variants", h.CreateVariant)
}

// List handles GET /products
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	companyIDStr := r.URL.Query().Get("company_id")
	if companyIDStr == "" {
		http.Error(w, "company_id is required", http.StatusBadRequest)
		return
	}

	companyID, err := uuid.Parse(companyIDStr)
	if err != nil {
		http.Error(w, "Invalid company_id", http.StatusBadRequest)
		return
	}

	limit := 50
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	products, err := h.service.List(r.Context(), companyID, limit, offset)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list products")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// GetByID handles GET /products/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetByID(r.Context(), id)
	if err == ErrNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get product")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Create handles POST /products
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		req.Status = "active"
	}

	product, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create product")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

// ListVariants handles GET /products/{id}/variants
func (h *Handler) ListVariants(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	variants, err := h.service.ListVariants(r.Context(), productID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list variants")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(variants)
}

// CreateVariant handles POST /products/{id}/variants
func (h *Handler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	productID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var req CreateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		req.Status = "active"
	}

	variant, err := h.service.CreateVariant(r.Context(), productID, req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create variant")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(variant)
}
