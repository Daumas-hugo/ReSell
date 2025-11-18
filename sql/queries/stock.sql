-- name: GetStockLevel :one
SELECT * FROM stock_levels
WHERE variant_id = $1 AND location_id = $2
LIMIT 1;

-- name: ListStockLevelsByVariant :many
SELECT * FROM stock_levels
WHERE variant_id = $1;

-- name: CreateStockLevel :one
INSERT INTO stock_levels (
    variant_id, location_id, quantity, reserved_quantity
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateStockLevel :one
UPDATE stock_levels
SET quantity = $2,
    reserved_quantity = $3,
    updated_at = NOW()
WHERE variant_id = $1 AND location_id = $4
RETURNING *;

-- name: ReserveStock :one
UPDATE stock_levels
SET reserved_quantity = reserved_quantity + $2,
    updated_at = NOW()
WHERE variant_id = $1 AND location_id = $3
RETURNING *;

-- name: ReleaseStock :one
UPDATE stock_levels
SET reserved_quantity = reserved_quantity - $2,
    updated_at = NOW()
WHERE variant_id = $1 AND location_id = $3
RETURNING *;

-- name: GetStockLocation :one
SELECT * FROM stock_locations
WHERE id = $1 LIMIT 1;

-- name: ListStockLocations :many
SELECT * FROM stock_locations
WHERE company_id = $1 AND status = $2
ORDER BY name;

-- name: CreateStockLocation :one
INSERT INTO stock_locations (
    company_id, name, code, type, status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;
