# Quick Start Guide

Get ReSell up and running in 5 minutes!

> **Windows Users**: This guide uses Unix-style commands. For Windows-specific instructions with PowerShell commands, see the [Windows Installation Guide](WINDOWS_SETUP.md).

## Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for local development)
- make (optional, but recommended)

## Step 1: Clone and Setup

```bash
git clone https://github.com/Daumas-hugo/ReSell.git
cd ReSell
cp .env.example .env
```

## Step 2: Start Infrastructure

Start PostgreSQL, Redis, and Keycloak using Docker Compose:

```bash
make docker-up
# or
docker-compose up -d
```

Wait for services to be healthy (about 60 seconds for Keycloak):

```bash
docker-compose ps
```

## Step 3: Run Migrations

Apply database migrations:

```bash
make migrate-up
# or
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/resell?sslmode=disable" up
```

## Step 4: Start the API

```bash
make run
# or
go run cmd/api/main.go
```

The API will start on `http://localhost:8080`

## Step 5: Test the API

### Health Check

```bash
curl http://localhost:8080/health
```

### Create a Company

```bash
curl -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "L'\''Heureux",
    "legal_name": "L'\''Heureux SAS",
    "siret": "12345678901234",
    "vat_number": "FR12345678901",
    "status": "active"
  }'
```

Save the returned company ID for next steps.

### Create a Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "YOUR_COMPANY_ID",
    "name": "Medical Device XYZ",
    "description": "High-quality medical device",
    "category": "medical_equipment",
    "weight_kg": 2.5,
    "status": "active"
  }'
```

Save the returned product ID.

### Create a Product Variant

```bash
curl -X POST http://localhost:8080/api/v1/products/YOUR_PRODUCT_ID/variants \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "MED-XYZ-001",
    "ean": "1234567890123",
    "base_price_ht": 150.00,
    "status": "active"
  }'
```

### List Companies

```bash
curl http://localhost:8080/api/v1/companies
```

### List Products

```bash
curl "http://localhost:8080/api/v1/products?company_id=YOUR_COMPANY_ID"
```

## Access Services

### PostgreSQL

```bash
psql -h localhost -p 5432 -U postgres -d resell
# Password: postgres
```

Or use a GUI like pgAdmin or DBeaver:
- Host: localhost
- Port: 5432
- Database: resell
- User: postgres
- Password: postgres

### Redis

```bash
redis-cli -h localhost -p 6379
```

### Keycloak Admin Console

Open http://localhost:8180/admin
- Username: admin
- Password: admin

## Stopping Services

```bash
make docker-down
# or
docker-compose down
```

## Development Workflow

### Making Code Changes

1. Edit code in `internal/` or `cmd/`
2. Build: `make build`
3. Test: `make test`
4. Run: `make run`

### Creating a Migration

```bash
make migrate-create NAME=add_new_feature
```

This creates two files in `migrations/`:
- `NNNNNN_add_new_feature.up.sql`
- `NNNNNN_add_new_feature.down.sql`

Edit the files and then:

```bash
make migrate-up
```

### Generating SQLC Code

After modifying queries in `sql/queries/`:

```bash
make sqlc
```

## Troubleshooting

### Port Already in Use

If ports 5432, 6379, or 8180 are already in use:

1. Edit `docker-compose.yml` to change ports
2. Update `.env` file with new ports
3. Restart: `make docker-down && make docker-up`

### Database Connection Failed

Check if PostgreSQL is running:

```bash
docker-compose ps postgres
```

Check logs:

```bash
docker-compose logs postgres
```

### Migration Errors

Reset database (WARNING: deletes all data):

```bash
make migrate-down
make migrate-up
```

Or force a specific version:

```bash
make migrate-force VERSION=1
```

## Next Steps

- Read the [API Documentation](API.md)
- Review the [Architecture Guide](ARCHITECTURE.md)
- Check the [README](../README.md) for more details
- Start building your e-commerce features!

## Example: Complete Order Flow

Here's a complete example of creating an order:

```bash
# 1. Create company
COMPANY_ID=$(curl -s -X POST http://localhost:8080/api/v1/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Company",
    "legal_name": "Test Company SAS",
    "status": "active"
  }' | jq -r '.id')

echo "Company ID: $COMPANY_ID"

# 2. Create customer (you'll need to add customer endpoints or insert via SQL)
# For now, insert a customer via SQL:
# psql -h localhost -U postgres -d resell -c "
#   INSERT INTO customers (company_id, type, email, first_name, last_name, status)
#   VALUES ('$COMPANY_ID', 'B2C', 'customer@example.com', 'John', 'Doe', 'active');
# "

# 3. Create product
PRODUCT_ID=$(curl -s -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d "{
    \"company_id\": \"$COMPANY_ID\",
    \"name\": \"Test Product\",
    \"status\": \"active\"
  }" | jq -r '.id')

echo "Product ID: $PRODUCT_ID"

# 4. Create variant
VARIANT_ID=$(curl -s -X POST http://localhost:8080/api/v1/products/$PRODUCT_ID/variants \
  -H "Content-Type: application/json" \
  -d '{
    "sku": "TEST-001",
    "base_price_ht": 100.00,
    "status": "active"
  }' | jq -r '.id')

echo "Variant ID: $VARIANT_ID"

# 5. Create order (need customer_id from step 2)
# 6. Add items to order
# etc...
```

## Tips

- Use `make help` to see all available commands
- Check logs with `docker-compose logs -f api`
- Use Postman or Insomnia for API testing
- Enable debug logging by setting `RESELL_LOG_LEVEL=debug`
