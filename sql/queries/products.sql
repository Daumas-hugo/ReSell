-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
WHERE company_id = $1 AND status = $2
ORDER BY name
LIMIT $3 OFFSET $4;

-- name: CreateProduct :one
INSERT INTO products (
    company_id, name, description, category,
    weight_kg, length_cm, width_cm, height_cm, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET name = $2,
    description = $3,
    category = $4,
    weight_kg = $5,
    length_cm = $6,
    width_cm = $7,
    height_cm = $8,
    status = $9,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: GetProductVariant :one
SELECT * FROM product_variants
WHERE id = $1 LIMIT 1;

-- name: GetProductVariantBySKU :one
SELECT * FROM product_variants
WHERE sku = $1 LIMIT 1;

-- name: ListProductVariants :many
SELECT * FROM product_variants
WHERE product_id = $1
ORDER BY created_at;

-- name: CreateProductVariant :one
INSERT INTO product_variants (
    product_id, sku, ean, manufacturer_code, base_price_ht, status
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: UpdateProductVariant :one
UPDATE product_variants
SET sku = $2,
    ean = $3,
    manufacturer_code = $4,
    base_price_ht = $5,
    status = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteProductVariant :exec
DELETE FROM product_variants
WHERE id = $1;
