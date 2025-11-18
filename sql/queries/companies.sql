-- name: GetCompany :one
SELECT * FROM companies
WHERE id = $1 LIMIT 1;

-- name: ListCompanies :many
SELECT * FROM companies
WHERE status = $1
ORDER BY name;

-- name: CreateCompany :one
INSERT INTO companies (
    name, legal_name, siret, vat_number, status
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateCompany :one
UPDATE companies
SET name = $2,
    legal_name = $3,
    siret = $4,
    vat_number = $5,
    status = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE id = $1;
