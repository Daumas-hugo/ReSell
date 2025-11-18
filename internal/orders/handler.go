package orders

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Handler handles HTTP requests for orders
type Handler struct {
	service *Service
	logger  zerolog.Logger
}

// NewHandler creates a new order handler
func NewHandler(service *Service, logger zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// RegisterRoutes registers order routes
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/{id}", h.GetByID)
	r.Get("/{id}/items", h.GetItems)
	r.Post("/{id}/items", h.AddItem)
}

// List handles GET /orders
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	companyIDStr := r.URL.Query().Get("company_id")
	customerIDStr := r.URL.Query().Get("customer_id")

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

	var orders []*Order
	var err error

	if customerIDStr != "" {
		customerID, parseErr := uuid.Parse(customerIDStr)
		if parseErr != nil {
			http.Error(w, "Invalid customer_id", http.StatusBadRequest)
			return
		}
		orders, err = h.service.ListByCustomer(r.Context(), customerID, limit, offset)
	} else if companyIDStr != "" {
		companyID, parseErr := uuid.Parse(companyIDStr)
		if parseErr != nil {
			http.Error(w, "Invalid company_id", http.StatusBadRequest)
			return
		}
		orders, err = h.service.ListByCompany(r.Context(), companyID, limit, offset)
	} else {
		http.Error(w, "company_id or customer_id is required", http.StatusBadRequest)
		return
	}

	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to list orders")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

// GetByID handles GET /orders/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetByID(r.Context(), id)
	if err == ErrNotFound {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get order")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// Create handles POST /orders
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.Create(r.Context(), req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create order")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetItems handles GET /orders/{id}/items
func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	items, err := h.service.GetItems(r.Context(), orderID)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to get order items")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// AddItem handles POST /orders/{id}/items
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	var req CreateOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	item, err := h.service.AddItem(r.Context(), orderID, req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to add order item")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}
