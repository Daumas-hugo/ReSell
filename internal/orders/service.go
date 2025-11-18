package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound = errors.New("order not found")
)

// Order represents an order
type Order struct {
	ID                 uuid.UUID
	CompanyID          uuid.UUID
	CustomerID         uuid.UUID
	OrderNumber        string
	Status             string
	SalesChannelID     *uuid.UUID
	SubtotalHT         decimal.Decimal
	DiscountHT         decimal.Decimal
	TaxAmount          decimal.Decimal
	TotalTTC           decimal.Decimal
	Currency           string
	BillingAddressID   *uuid.UUID
	ShippingAddressID  *uuid.UUID
	Notes              *string
	CreatedAt          time.Time
}

// OrderItem represents an order item
type OrderItem struct {
	ID           uuid.UUID
	OrderID      uuid.UUID
	VariantID    uuid.UUID
	SKU          string
	ProductName  string
	VariantName  *string
	Quantity     int
	UnitPriceHT  decimal.Decimal
	UnitPriceTTC decimal.Decimal
	DiscountHT   decimal.Decimal
	TaxRate      decimal.Decimal
	TaxAmount    decimal.Decimal
	TotalHT      decimal.Decimal
	TotalTTC     decimal.Decimal
}

// CreateOrderRequest is the request to create an order
type CreateOrderRequest struct {
	CompanyID         uuid.UUID       `json:"company_id"`
	CustomerID        uuid.UUID       `json:"customer_id"`
	SalesChannelID    *uuid.UUID      `json:"sales_channel_id,omitempty"`
	BillingAddressID  *uuid.UUID      `json:"billing_address_id,omitempty"`
	ShippingAddressID *uuid.UUID      `json:"shipping_address_id,omitempty"`
	Notes             *string         `json:"notes,omitempty"`
}

// CreateOrderItemRequest is the request to add an item to an order
type CreateOrderItemRequest struct {
	VariantID    uuid.UUID       `json:"variant_id"`
	SKU          string          `json:"sku"`
	ProductName  string          `json:"product_name"`
	VariantName  *string         `json:"variant_name,omitempty"`
	Quantity     int             `json:"quantity"`
	UnitPriceHT  decimal.Decimal `json:"unit_price_ht"`
	UnitPriceTTC decimal.Decimal `json:"unit_price_ttc"`
	DiscountHT   decimal.Decimal `json:"discount_ht"`
	TaxRate      decimal.Decimal `json:"tax_rate"`
	TaxAmount    decimal.Decimal `json:"tax_amount"`
	TotalHT      decimal.Decimal `json:"total_ht"`
	TotalTTC     decimal.Decimal `json:"total_ttc"`
}

// Service handles order business logic
type Service struct {
	db *sql.DB
}

// NewService creates a new order service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetByID retrieves an order by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Order, error) {
	query := `
		SELECT id, company_id, customer_id, order_number, status, sales_channel_id,
		       subtotal_ht, discount_ht, tax_amount, total_ttc, currency,
		       billing_address_id, shipping_address_id, notes, created_at
		FROM orders
		WHERE id = $1
	`
	
	var order Order
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.CompanyID,
		&order.CustomerID,
		&order.OrderNumber,
		&order.Status,
		&order.SalesChannelID,
		&order.SubtotalHT,
		&order.DiscountHT,
		&order.TaxAmount,
		&order.TotalTTC,
		&order.Currency,
		&order.BillingAddressID,
		&order.ShippingAddressID,
		&order.Notes,
		&order.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

// ListByCompany retrieves orders for a company
func (s *Service) ListByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*Order, error) {
	query := `
		SELECT id, company_id, customer_id, order_number, status, sales_channel_id,
		       subtotal_ht, discount_ht, tax_amount, total_ttc, currency,
		       billing_address_id, shipping_address_id, notes, created_at
		FROM orders
		WHERE company_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := s.db.QueryContext(ctx, query, companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []*Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(
			&order.ID,
			&order.CompanyID,
			&order.CustomerID,
			&order.OrderNumber,
			&order.Status,
			&order.SalesChannelID,
			&order.SubtotalHT,
			&order.DiscountHT,
			&order.TaxAmount,
			&order.TotalTTC,
			&order.Currency,
			&order.BillingAddressID,
			&order.ShippingAddressID,
			&order.Notes,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	
	return orders, nil
}

// ListByCustomer retrieves orders for a customer
func (s *Service) ListByCustomer(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*Order, error) {
	query := `
		SELECT id, company_id, customer_id, order_number, status, sales_channel_id,
		       subtotal_ht, discount_ht, tax_amount, total_ttc, currency,
		       billing_address_id, shipping_address_id, notes, created_at
		FROM orders
		WHERE customer_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := s.db.QueryContext(ctx, query, customerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []*Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(
			&order.ID,
			&order.CompanyID,
			&order.CustomerID,
			&order.OrderNumber,
			&order.Status,
			&order.SalesChannelID,
			&order.SubtotalHT,
			&order.DiscountHT,
			&order.TaxAmount,
			&order.TotalTTC,
			&order.Currency,
			&order.BillingAddressID,
			&order.ShippingAddressID,
			&order.Notes,
			&order.CreatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}
	
	return orders, nil
}

// Create creates a new order
func (s *Service) Create(ctx context.Context, req CreateOrderRequest) (*Order, error) {
	// Generate order number
	orderNumber := fmt.Sprintf("ORD-%d", time.Now().Unix())
	
	query := `
		INSERT INTO orders (company_id, customer_id, order_number, status,
		                    sales_channel_id, billing_address_id, shipping_address_id, notes)
		VALUES ($1, $2, $3, 'draft', $4, $5, $6, $7)
		RETURNING id, company_id, customer_id, order_number, status, sales_channel_id,
		          subtotal_ht, discount_ht, tax_amount, total_ttc, currency,
		          billing_address_id, shipping_address_id, notes, created_at
	`
	
	var order Order
	err := s.db.QueryRowContext(ctx, query,
		req.CompanyID,
		req.CustomerID,
		orderNumber,
		req.SalesChannelID,
		req.BillingAddressID,
		req.ShippingAddressID,
		req.Notes,
	).Scan(
		&order.ID,
		&order.CompanyID,
		&order.CustomerID,
		&order.OrderNumber,
		&order.Status,
		&order.SalesChannelID,
		&order.SubtotalHT,
		&order.DiscountHT,
		&order.TaxAmount,
		&order.TotalTTC,
		&order.Currency,
		&order.BillingAddressID,
		&order.ShippingAddressID,
		&order.Notes,
		&order.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &order, nil
}

// GetItems retrieves items for an order
func (s *Service) GetItems(ctx context.Context, orderID uuid.UUID) ([]*OrderItem, error) {
	query := `
		SELECT id, order_id, variant_id, sku, product_name, variant_name,
		       quantity, unit_price_ht, unit_price_ttc, discount_ht,
		       tax_rate, tax_amount, total_ht, total_ttc
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at
	`
	
	rows, err := s.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var items []*OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.VariantID,
			&item.SKU,
			&item.ProductName,
			&item.VariantName,
			&item.Quantity,
			&item.UnitPriceHT,
			&item.UnitPriceTTC,
			&item.DiscountHT,
			&item.TaxRate,
			&item.TaxAmount,
			&item.TotalHT,
			&item.TotalTTC,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	
	return items, nil
}

// AddItem adds an item to an order
func (s *Service) AddItem(ctx context.Context, orderID uuid.UUID, req CreateOrderItemRequest) (*OrderItem, error) {
	query := `
		INSERT INTO order_items (
			order_id, variant_id, sku, product_name, variant_name,
			quantity, unit_price_ht, unit_price_ttc, discount_ht,
			tax_rate, tax_amount, total_ht, total_ttc
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, order_id, variant_id, sku, product_name, variant_name,
		          quantity, unit_price_ht, unit_price_ttc, discount_ht,
		          tax_rate, tax_amount, total_ht, total_ttc
	`
	
	var item OrderItem
	err := s.db.QueryRowContext(ctx, query,
		orderID,
		req.VariantID,
		req.SKU,
		req.ProductName,
		req.VariantName,
		req.Quantity,
		req.UnitPriceHT,
		req.UnitPriceTTC,
		req.DiscountHT,
		req.TaxRate,
		req.TaxAmount,
		req.TotalHT,
		req.TotalTTC,
	).Scan(
		&item.ID,
		&item.OrderID,
		&item.VariantID,
		&item.SKU,
		&item.ProductName,
		&item.VariantName,
		&item.Quantity,
		&item.UnitPriceHT,
		&item.UnitPriceTTC,
		&item.DiscountHT,
		&item.TaxRate,
		&item.TaxAmount,
		&item.TotalHT,
		&item.TotalTTC,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &item, nil
}
