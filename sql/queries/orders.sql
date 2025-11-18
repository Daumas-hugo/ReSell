-- name: GetOrder :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: GetOrderByNumber :one
SELECT * FROM orders
WHERE order_number = $1 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
WHERE company_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListCustomerOrders :many
SELECT * FROM orders
WHERE customer_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateOrder :one
INSERT INTO orders (
    company_id, customer_id, order_number, status,
    sales_channel_id, subtotal_ht, discount_ht, tax_amount,
    total_ttc, currency, billing_address_id, shipping_address_id, notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrderTotals :one
UPDATE orders
SET subtotal_ht = $2,
    discount_ht = $3,
    tax_amount = $4,
    total_ttc = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1;

-- name: GetOrderItem :one
SELECT * FROM order_items
WHERE id = $1 LIMIT 1;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1
ORDER BY created_at;

-- name: CreateOrderItem :one
INSERT INTO order_items (
    order_id, variant_id, sku, product_name, variant_name,
    quantity, unit_price_ht, unit_price_ttc, discount_ht,
    tax_rate, tax_amount, total_ht, total_ttc
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE id = $1;
