# ReSell - E-Commerce Backend Monolith

A modern, modular e-commerce backend built with Go, designed for B2B and B2C operations with multi-company support.

## Features

- **Multi-Company Support**: Manage multiple brands/companies in a single platform
- **Unified B2B/B2C Portal**: Single portal supporting both business and consumer customers
- **Shopify-like Product Model**: Products with variants, options, and media
- **Flexible Pricing Engine**: Calculate prices based on customer type, sales channel, promotions, and coupons
- **Pluggable Stock Management**: Connector-based architecture for WMS/ERP integration
- **Keycloak Authentication**: OIDC-based authentication and authorization

## Architecture

### Tech Stack

- **Go 1.24+**: Modern, efficient backend
- **PostgreSQL 16**: Primary database
- **Redis 7**: Caching and messaging
- **Keycloak 24+**: Identity and access management
- **Chi v5**: HTTP router
- **SQLC**: Type-safe SQL code generation
- **Zerolog**: Structured logging
- **Viper**: Configuration management

### Project Structure

```
.
├── cmd/
│   └── api/                 # Main application entry point
├── internal/
│   ├── auth/               # Authentication & authorization
│   ├── companies/          # Company management
│   ├── customers/          # Customer management (B2B/B2C)
│   ├── products/           # Product & variant management
│   ├── pricing/            # Pricing calculation engine
│   ├── stock/              # Stock management connectors
│   ├── orders/             # Order processing
│   ├── config/             # Configuration
│   └── database/           # Database connection
├── pkg/
│   ├── logger/             # Logging utilities
│   └── middleware/         # HTTP middleware
├── migrations/             # Database migrations
├── sql/
│   ├── queries/            # SQLC queries
│   └── schema/             # Database schema
└── deployments/
    └── docker/             # Docker configurations
```

## Getting Started

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- make (optional, for convenience)

> **Windows Users**: See the [Windows Installation Guide](docs/WINDOWS_SETUP.md) for detailed Windows-specific setup instructions.

### Installation

#### Linux / macOS

1. Clone the repository:
```bash
git clone https://github.com/Daumas-hugo/ReSell.git
cd ReSell
```

2. Install dependencies:
```bash
make deps
```

3. Install development tools:
```bash
make install-tools
```

### Development Environment

1. Start infrastructure services (PostgreSQL, Redis, Keycloak):
```bash
make docker-up
```

2. Run database migrations:
```bash
make migrate-up
```

3. (Optional) Generate SQLC code if queries were modified:
```bash
make sqlc
```

4. Start the API server:
```bash
make run
```

The API will be available at `http://localhost:8080`

### Environment Variables

Copy `.env.example` to `.env` and adjust values as needed:

```bash
cp .env.example .env
```

Key configuration options:
- `RESELL_SERVER_PORT`: API server port (default: 8080)
- `RESELL_DATABASE_*`: PostgreSQL connection settings
- `RESELL_REDIS_*`: Redis connection settings
- `RESELL_KEYCLOAK_*`: Keycloak integration settings

## API Endpoints

### Health Check
```
GET /health
```

### Companies
```
GET    /api/v1/companies
POST   /api/v1/companies
GET    /api/v1/companies/:id
PUT    /api/v1/companies/:id
DELETE /api/v1/companies/:id
```

### Customers
```
GET    /api/v1/customers
POST   /api/v1/customers
GET    /api/v1/customers/:id
PUT    /api/v1/customers/:id
DELETE /api/v1/customers/:id
```

### Products
```
GET    /api/v1/products
POST   /api/v1/products
GET    /api/v1/products/:id
PUT    /api/v1/products/:id
DELETE /api/v1/products/:id
GET    /api/v1/products/:id/variants
POST   /api/v1/products/:id/variants
```

### Orders
```
GET    /api/v1/orders
POST   /api/v1/orders
GET    /api/v1/orders/:id
PUT    /api/v1/orders/:id
DELETE /api/v1/orders/:id
```

### Customer Portal
```
GET /api/v1/me              # Current user info
GET /api/v1/me/orders       # User's orders
GET /api/v1/me/companies    # User's companies
```

## Database Schema

The database schema includes:

- **companies**: Multi-tenancy support
- **customers**: B2B/B2C customer profiles
- **products**: Product catalog
- **product_variants**: Product SKUs with pricing
- **product_options**: Product variations (size, color, etc.)
- **orders**: Customer orders
- **pricing_rules**: Dynamic pricing rules
- **promotions**: Marketing promotions and coupons
- **stock_levels**: Inventory management
- **tax_rates**: Tax calculation

## Modules

### Pricing Engine

The pricing engine calculates prices based on multiple factors:
- Base variant price
- Customer type (B2B/B2C)
- Sales channel (web, portal, mobile, etc.)
- Volume discounts
- Promotional offers
- Coupon codes
- Tax rates by country

Example usage:
```go
priceReq := pricing.PriceRequest{
    VariantID:    variantID,
    CompanyID:    companyID,
    CustomerType: pricing.CustomerTypeB2B,
    SalesChannel: pricing.SalesChannelWeb,
    Quantity:     10,
    Country:      "FR",
}

result, err := pricingEngine.CalculatePrice(ctx, priceReq)
```

### Stock Management

The stock module uses a connector pattern for flexibility:
- **LocalConnector**: Database-backed stock management
- **MockConnector**: In-memory stock for testing
- **WMSConnector**: (Future) Connect to warehouse management systems
- **ERPConnector**: (Future) Connect to ERP systems

Example usage:
```go
stock, err := stockManager.GetAvailableStock(ctx, variantID, companyID, nil)
if err != nil {
    // Handle error
}

err = stockManager.ReserveStock(ctx, variantID, locationID, quantity)
```

## Testing

Run tests:
```bash
make test
```

Run tests with coverage:
```bash
make test-coverage
```

## Building for Production

Build the binary:
```bash
make build
```

Build Docker image:
```bash
make docker-build
```

## Database Migrations

Create a new migration:
```bash
make migrate-create NAME=add_new_table
```

Apply migrations:
```bash
make migrate-up
```

Rollback last migration:
```bash
make migrate-down
```

## Contributing

This project is designed to be modular and extensible. When contributing:

1. Follow the existing code structure
2. Write tests for new features
3. Update documentation
4. Ensure all tests pass
5. Run `make fmt` and `make lint`

## License

This project is designed to be open-source ready. License TBD.

## Roadmap

- [ ] Complete API endpoint implementations
- [ ] Authentication/authorization middleware with Keycloak
- [ ] Advanced pricing rules engine
- [ ] Real WMS/ERP connectors
- [ ] Order fulfillment workflow
- [ ] Invoice generation
- [ ] Customer portal frontend
- [ ] Admin panel
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Performance optimization
- [ ] Comprehensive test coverage
