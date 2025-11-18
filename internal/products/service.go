package products

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

var (
	ErrNotFound = errors.New("product not found")
)

// Product represents a product
type Product struct {
	ID          uuid.UUID
	CompanyID   uuid.UUID
	Name        string
	Description *string
	Category    *string
	WeightKg    *decimal.Decimal
	LengthCm    *decimal.Decimal
	WidthCm     *decimal.Decimal
	HeightCm    *decimal.Decimal
	Status      string
}

// ProductVariant represents a product variant
type ProductVariant struct {
	ID               uuid.UUID
	ProductID        uuid.UUID
	SKU              string
	EAN              *string
	ManufacturerCode *string
	BasePriceHT      decimal.Decimal
	Status           string
}

// CreateProductRequest is the request to create a product
type CreateProductRequest struct {
	CompanyID   uuid.UUID        `json:"company_id"`
	Name        string           `json:"name"`
	Description *string          `json:"description,omitempty"`
	Category    *string          `json:"category,omitempty"`
	WeightKg    *decimal.Decimal `json:"weight_kg,omitempty"`
	LengthCm    *decimal.Decimal `json:"length_cm,omitempty"`
	WidthCm     *decimal.Decimal `json:"width_cm,omitempty"`
	HeightCm    *decimal.Decimal `json:"height_cm,omitempty"`
	Status      string           `json:"status"`
}

// CreateVariantRequest is the request to create a variant
type CreateVariantRequest struct {
	SKU              string          `json:"sku"`
	EAN              *string         `json:"ean,omitempty"`
	ManufacturerCode *string         `json:"manufacturer_code,omitempty"`
	BasePriceHT      decimal.Decimal `json:"base_price_ht"`
	Status           string          `json:"status"`
}

// Service handles product business logic
type Service struct {
	db *sql.DB
}

// NewService creates a new product service
func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

// GetByID retrieves a product by ID
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	query := `
		SELECT id, company_id, name, description, category, 
		       weight_kg, length_cm, width_cm, height_cm, status
		FROM products
		WHERE id = $1
	`
	
	var product Product
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.CompanyID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.WeightKg,
		&product.LengthCm,
		&product.WidthCm,
		&product.HeightCm,
		&product.Status,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	if err != nil {
		return nil, err
	}
	
	return &product, nil
}

// List retrieves products for a company
func (s *Service) List(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*Product, error) {
	query := `
		SELECT id, company_id, name, description, category,
		       weight_kg, length_cm, width_cm, height_cm, status
		FROM products
		WHERE company_id = $1 AND status = 'active'
		ORDER BY name
		LIMIT $2 OFFSET $3
	`
	
	rows, err := s.db.QueryContext(ctx, query, companyID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var products []*Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ID,
			&product.CompanyID,
			&product.Name,
			&product.Description,
			&product.Category,
			&product.WeightKg,
			&product.LengthCm,
			&product.WidthCm,
			&product.HeightCm,
			&product.Status,
		); err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	
	return products, nil
}

// Create creates a new product
func (s *Service) Create(ctx context.Context, req CreateProductRequest) (*Product, error) {
	query := `
		INSERT INTO products (company_id, name, description, category,
		                      weight_kg, length_cm, width_cm, height_cm, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, company_id, name, description, category,
		          weight_kg, length_cm, width_cm, height_cm, status
	`
	
	var product Product
	err := s.db.QueryRowContext(ctx, query,
		req.CompanyID,
		req.Name,
		req.Description,
		req.Category,
		req.WeightKg,
		req.LengthCm,
		req.WidthCm,
		req.HeightCm,
		req.Status,
	).Scan(
		&product.ID,
		&product.CompanyID,
		&product.Name,
		&product.Description,
		&product.Category,
		&product.WeightKg,
		&product.LengthCm,
		&product.WidthCm,
		&product.HeightCm,
		&product.Status,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &product, nil
}

// GetVariant retrieves a variant by ID
func (s *Service) GetVariant(ctx context.Context, id uuid.UUID) (*ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, ean, manufacturer_code, base_price_ht, status
		FROM product_variants
		WHERE id = $1
	`
	
	var variant ProductVariant
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.EAN,
		&variant.ManufacturerCode,
		&variant.BasePriceHT,
		&variant.Status,
	)
	
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	
	if err != nil {
		return nil, err
	}
	
	return &variant, nil
}

// ListVariants retrieves variants for a product
func (s *Service) ListVariants(ctx context.Context, productID uuid.UUID) ([]*ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, ean, manufacturer_code, base_price_ht, status
		FROM product_variants
		WHERE product_id = $1
		ORDER BY created_at
	`
	
	rows, err := s.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var variants []*ProductVariant
	for rows.Next() {
		var variant ProductVariant
		if err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.SKU,
			&variant.EAN,
			&variant.ManufacturerCode,
			&variant.BasePriceHT,
			&variant.Status,
		); err != nil {
			return nil, err
		}
		variants = append(variants, &variant)
	}
	
	return variants, nil
}

// CreateVariant creates a new variant
func (s *Service) CreateVariant(ctx context.Context, productID uuid.UUID, req CreateVariantRequest) (*ProductVariant, error) {
	query := `
		INSERT INTO product_variants (product_id, sku, ean, manufacturer_code, base_price_ht, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, product_id, sku, ean, manufacturer_code, base_price_ht, status
	`
	
	var variant ProductVariant
	err := s.db.QueryRowContext(ctx, query,
		productID,
		req.SKU,
		req.EAN,
		req.ManufacturerCode,
		req.BasePriceHT,
		req.Status,
	).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.EAN,
		&variant.ManufacturerCode,
		&variant.BasePriceHT,
		&variant.Status,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &variant, nil
}
