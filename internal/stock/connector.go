package stock

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrStockNotAvailable = errors.New("stock not available")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrLocationNotFound  = errors.New("location not found")
)

// StockLevel represents available stock information
type StockLevel struct {
	VariantID        uuid.UUID
	LocationID       uuid.UUID
	Quantity         int
	ReservedQuantity int
	AvailableQty     int
}

// Connector defines the interface for stock management systems
type Connector interface {
	// GetAvailableStock returns available stock for a variant
	GetAvailableStock(ctx context.Context, variantID, companyID uuid.UUID, locationID *uuid.UUID) (*StockLevel, error)
	
	// ReserveStock reserves stock for an order
	ReserveStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error
	
	// ReleaseStock releases previously reserved stock
	ReleaseStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error
	
	// UpdateStock updates stock levels (for inventory adjustments)
	UpdateStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error
}

// Manager manages stock operations with multiple connectors
type Manager struct {
	primary  Connector
	fallback Connector
}

// NewManager creates a new stock manager
func NewManager(primary Connector, fallback Connector) *Manager {
	return &Manager{
		primary:  primary,
		fallback: fallback,
	}
}

// GetAvailableStock gets available stock, trying primary then fallback
func (m *Manager) GetAvailableStock(ctx context.Context, variantID, companyID uuid.UUID, locationID *uuid.UUID) (*StockLevel, error) {
	stock, err := m.primary.GetAvailableStock(ctx, variantID, companyID, locationID)
	if err == nil {
		return stock, nil
	}

	if m.fallback != nil {
		return m.fallback.GetAvailableStock(ctx, variantID, companyID, locationID)
	}

	return nil, err
}

// ReserveStock reserves stock using primary connector
func (m *Manager) ReserveStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	return m.primary.ReserveStock(ctx, variantID, locationID, quantity)
}

// ReleaseStock releases stock using primary connector
func (m *Manager) ReleaseStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	return m.primary.ReleaseStock(ctx, variantID, locationID, quantity)
}

// UpdateStock updates stock levels
func (m *Manager) UpdateStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	return m.primary.UpdateStock(ctx, variantID, locationID, quantity)
}
