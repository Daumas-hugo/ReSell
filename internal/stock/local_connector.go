package stock

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

// LocalConnector implements Connector using local database
type LocalConnector struct {
	db *sql.DB
}

// NewLocalConnector creates a new local database connector
func NewLocalConnector(db *sql.DB) *LocalConnector {
	return &LocalConnector{db: db}
}

// GetAvailableStock gets stock from local database
func (c *LocalConnector) GetAvailableStock(ctx context.Context, variantID, companyID uuid.UUID, locationID *uuid.UUID) (*StockLevel, error) {
	var stock StockLevel
	var query string
	var args []interface{}

	if locationID != nil {
		query = `
			SELECT variant_id, location_id, quantity, reserved_quantity,
			       (quantity - reserved_quantity) as available_qty
			FROM stock_levels
			WHERE variant_id = $1 AND location_id = $2
		`
		args = []interface{}{variantID, *locationID}
	} else {
		// Sum across all locations
		query = `
			SELECT variant_id, location_id, 
			       SUM(quantity) as quantity, 
			       SUM(reserved_quantity) as reserved_quantity,
			       SUM(quantity - reserved_quantity) as available_qty
			FROM stock_levels sl
			JOIN stock_locations loc ON sl.location_id = loc.id
			WHERE sl.variant_id = $1 AND loc.company_id = $2
			GROUP BY variant_id, location_id
			LIMIT 1
		`
		args = []interface{}{variantID, companyID}
	}

	err := c.db.QueryRowContext(ctx, query, args...).Scan(
		&stock.VariantID,
		&stock.LocationID,
		&stock.Quantity,
		&stock.ReservedQuantity,
		&stock.AvailableQty,
	)

	if err == sql.ErrNoRows {
		return nil, ErrStockNotAvailable
	}

	if err != nil {
		return nil, err
	}

	return &stock, nil
}

// ReserveStock reserves stock in local database
func (c *LocalConnector) ReserveStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	query := `
		UPDATE stock_levels
		SET reserved_quantity = reserved_quantity + $3,
		    updated_at = NOW()
		WHERE variant_id = $1 AND location_id = $2
		  AND (quantity - reserved_quantity) >= $3
	`

	result, err := c.db.ExecContext(ctx, query, variantID, locationID, quantity)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrInsufficientStock
	}

	return nil
}

// ReleaseStock releases reserved stock in local database
func (c *LocalConnector) ReleaseStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	query := `
		UPDATE stock_levels
		SET reserved_quantity = reserved_quantity - $3,
		    updated_at = NOW()
		WHERE variant_id = $1 AND location_id = $2
		  AND reserved_quantity >= $3
	`

	result, err := c.db.ExecContext(ctx, query, variantID, locationID, quantity)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrInsufficientStock
	}

	return nil
}

// UpdateStock updates stock quantity in local database
func (c *LocalConnector) UpdateStock(ctx context.Context, variantID, locationID uuid.UUID, quantity int) error {
	query := `
		UPDATE stock_levels
		SET quantity = $3,
		    updated_at = NOW()
		WHERE variant_id = $1 AND location_id = $2
	`

	_, err := c.db.ExecContext(ctx, query, variantID, locationID, quantity)
	return err
}
