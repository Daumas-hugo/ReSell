package stock

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

// MockConnector implements Connector for testing/development
type MockConnector struct {
	mu    sync.RWMutex
	stock map[string]*StockLevel
}

// NewMockConnector creates a new mock connector
func NewMockConnector() *MockConnector {
	return &MockConnector{
		stock: make(map[string]*StockLevel),
	}
}

// SetStock sets stock for testing
func (c *MockConnector) SetStock(variantID, locationID uuid.UUID, quantity int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := variantID.String() + ":" + locationID.String()
	c.stock[key] = &StockLevel{
		VariantID:        variantID,
		LocationID:       locationID,
		Quantity:         quantity,
		ReservedQuantity: 0,
		AvailableQty:     quantity,
	}
}

// GetAvailableStock returns mock stock
func (c *MockConnector) GetAvailableStock(ctx context.Context, variantID, companyID uuid.UUID, locationID *uuid.UUID) (*StockLevel, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if locationID == nil {
		// Return first match for variant
		for _, stock := range c.stock {
			if stock.VariantID == variantID {
				return stock, nil
			}
		}
		return nil, ErrStockNotAvailable
	}

	key := variantID.String() + ":" + locationID.String()
	stock, ok := c.stock[key]
	if !ok {
		return nil, ErrStockNotAvailable
	}

	return stock, nil
}

// ReserveStock reserves mock stock
func (c *MockConnector) ReserveStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := variantID.String() + ":" + locationID.String()
	stock, ok := c.stock[key]
	if !ok {
		return ErrStockNotAvailable
	}

	if stock.AvailableQty < quantity {
		return ErrInsufficientStock
	}

	stock.ReservedQuantity += quantity
	stock.AvailableQty = stock.Quantity - stock.ReservedQuantity

	return nil
}

// ReleaseStock releases mock stock
func (c *MockConnector) ReleaseStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := variantID.String() + ":" + locationID.String()
	stock, ok := c.stock[key]
	if !ok {
		return ErrStockNotAvailable
	}

	if stock.ReservedQuantity < quantity {
		return ErrInsufficientStock
	}

	stock.ReservedQuantity -= quantity
	stock.AvailableQty = stock.Quantity - stock.ReservedQuantity

	return nil
}

// UpdateStock updates mock stock
func (c *MockConnector) UpdateStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := variantID.String() + ":" + locationID.String()
	stock, ok := c.stock[key]
	if !ok {
		stock = &StockLevel{
			VariantID:  variantID,
			LocationID: locationID,
		}
		c.stock[key] = stock
	}

	stock.Quantity = quantity
	stock.AvailableQty = stock.Quantity - stock.ReservedQuantity

	return nil
}
