-- name: GetCustomer :one
SELECT * FROM customers
WHERE id = $1 LIMIT 1;

-- name: GetCustomerByEmail :one
SELECT * FROM customers
WHERE company_id = $1 AND email = $2
LIMIT 1;

-- name: GetCustomerByKeycloakID :one
SELECT * FROM customers
WHERE keycloak_user_id = $1
LIMIT 1;

-- name: ListCustomers :many
SELECT * FROM customers
WHERE company_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateCustomer :one
INSERT INTO customers (
    company_id, type, email, first_name, last_name,
    organization_name, siret, vat_number, keycloak_user_id, status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateCustomer :one
UPDATE customers
SET first_name = $2,
    last_name = $3,
    organization_name = $4,
    siret = $5,
    vat_number = $6,
    status = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customers
WHERE id = $1;
