# ReSell API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Health Check

### Check API Health

```
GET /health
```

Returns `200 OK` if the API is running.

## Companies

### List Companies

```
GET /api/v1/companies
```

Returns all active companies.

### Get Company

```
GET /api/v1/companies/{id}
```

Returns a single company by ID.

### Create Company

```
POST /api/v1/companies
```

**Request Body:**
```json
{
  "name": "Medimat",
  "legal_name": "Medimat SARL",
  "siret": "98765432109876",
  "vat_number": "FR98765432109",
  "status": "active"
}
```

### Update Company

```
PUT /api/v1/companies/{id}
```

### Delete Company

```
DELETE /api/v1/companies/{id}
```

## Products

### List Products

```
GET /api/v1/products?company_id={uuid}&limit=50&offset=0
```

### Get Product

```
GET /api/v1/products/{id}
```

### Create Product

```
POST /api/v1/products
```

### List Product Variants

```
GET /api/v1/products/{id}/variants
```

### Create Product Variant

```
POST /api/v1/products/{id}/variants
```

## Orders

### List Orders

```
GET /api/v1/orders?company_id={uuid}
GET /api/v1/orders?customer_id={uuid}
```

### Get Order

```
GET /api/v1/orders/{id}
```

### Create Order

```
POST /api/v1/orders
```

### Get Order Items

```
GET /api/v1/orders/{id}/items
```

### Add Order Item

```
POST /api/v1/orders/{id}/items
```

## Customer Portal

### Get Current User

```
GET /api/v1/me
```

**Status:** `501 Not Implemented` (authentication not yet configured)

### Get User Orders

```
GET /api/v1/me/orders
```

**Status:** `501 Not Implemented`

### Get User Companies

```
GET /api/v1/me/companies
```

**Status:** `501 Not Implemented`
