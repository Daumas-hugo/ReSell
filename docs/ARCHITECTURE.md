# ReSell Architecture

## Overview

ReSell is a monolithic e-commerce backend designed to support both B2B and B2C operations with multi-company (multi-tenancy) capabilities. The architecture follows clean architecture principles with clear separation of concerns.

## Architecture Principles

1. **Modular Monolith**: Each business domain is isolated in its own package with clear boundaries
2. **Database-Centric**: PostgreSQL as the single source of truth
3. **Type Safety**: SQLC for type-safe database operations
4. **Extensibility**: Plugin architecture for external systems (stock, payments, etc.)
5. **Configuration-Driven**: Environment-based configuration for different deployments

## Technology Stack

### Core
- **Language**: Go 1.24.x
- **HTTP Router**: Chi v5
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **Auth**: Keycloak 24+

### Development
- **Migrations**: golang-migrate
- **SQL Generation**: SQLC v1.27+
- **Config**: Viper
- **Logging**: zerolog
- **Containerization**: Docker & Docker Compose

## Project Structure

```
/cmd
  /api                 # Main application entry point
/internal              # Private application code
  /companies           # Company management module
    handler.go         # HTTP handlers
    service.go         # Business logic
  /customers           # Customer management (B2B/B2C)
  /products            # Product catalog & variants
  /pricing             # Pricing calculation engine
  /stock               # Stock management with connectors
  /orders              # Order processing
  /auth                # Authentication & authorization
  /config              # Configuration management
  /database            # Database connection & utilities
/pkg                   # Public reusable packages
  /logger              # Logging utilities
  /middleware          # HTTP middleware
/migrations            # Database migrations
/sql
  /queries             # SQLC query definitions
  /schema              # Database schema
/docs                  # Documentation
```

## Domain Models

### 1. Companies (Multi-Tenancy)

Companies represent different brands/entities within the platform. All data is scoped by `company_id`.

**Key Features:**
- Company profile (name, legal info, VAT, SIRET)
- Status management
- Isolation of data per company

### 2. Customers (B2B/B2C)

Unified customer model supporting both consumer and business customers.

**B2C Customers:**
- Individual profiles
- Simple checkout flow
- Personal data

**B2B Customers:**
- Organization info (SIRET, VAT)
- Multiple user contacts
- Role-based access
- Negotiated pricing

**Key Features:**
- Single table for both types
- Keycloak integration for identity
- Multi-company association
- Address management

### 3. Products & Variants

Shopify-inspired product model with flexible variants.

**Product:**
- Base product information
- Belongs to a company
- Has options (size, color, etc.)
- Has media (images, documents)
- Physical attributes (dimensions, weight)

**Product Variant:**
- Unique SKU
- Specific option combination
- Base price (before calculations)
- Stock tracking
- EAN/manufacturer codes

**Key Features:**
- Everything sold is a variant, not a product
- Flexible option system
- Media management
- Category organization

### 4. Pricing Engine

Dynamic price calculation based on multiple factors.

**Input Parameters:**
- Variant ID
- Company ID
- Customer type (B2B/B2C)
- Sales channel
- Quantity
- Coupon codes
- Country (for tax)

**Calculation Steps:**
1. Start with variant base price
2. Apply sales channel rules
3. Apply customer type rules (B2B discounts)
4. Apply volume discounts
5. Apply promotions
6. Apply coupon codes
7. Calculate taxes

**Output:**
- Unit price HT (without tax)
- Unit price TTC (with tax)
- Total HT
- Total TTC
- Tax details
- Applied discounts details

### 5. Stock Management

Connector-based architecture for flexibility.

**Stock Connector Interface:**
```go
type Connector interface {
    GetAvailableStock(ctx, variantID, companyID, locationID) (*StockLevel, error)
    ReserveStock(ctx, variantID, locationID, quantity) error
    ReleaseStock(ctx, variantID, locationID, quantity) error
    UpdateStock(ctx, variantID, locationID, quantity) error
}
```

**Available Connectors:**
- **LocalConnector**: Database-backed stock (for MVP)
- **MockConnector**: In-memory for testing
- **Future**: WMSConnector, ERPConnector

**Stock Manager:**
- Primary connector (usually local or WMS)
- Fallback connector (for high availability)
- Automatic failover

### 6. Orders

Complete order management system.

**Order:**
- Belongs to company and customer
- Order number (unique)
- Status workflow (draft → pending_payment → paid → shipped → cancelled)
- Sales channel tracking
- Address references
- Price totals (HT, TTC)

**Order Items:**
- References variant (soft delete protection)
- Captures price at order time (no recalculation)
- Item-level discounts
- Tax details per item

**Key Features:**
- Variant-based (never direct products)
- Frozen pricing
- Multi-currency support
- Notes/comments

## Data Flow

### Typical Order Creation Flow

```
1. Customer adds variant to cart
   ↓
2. Pricing engine calculates price
   - Gets variant base price
   - Applies rules and discounts
   - Calculates tax
   ↓
3. Stock check
   - Queries stock connector
   - Verifies availability
   ↓
4. Order creation
   - Creates order (draft status)
   - Adds order items (with calculated prices)
   ↓
5. Stock reservation
   - Reserves stock for order items
   ↓
6. Payment processing (external)
   ↓
7. Order confirmation
   - Status → pending_payment → paid
   ↓
8. Fulfillment (future)
```

## Database Design

### Key Design Decisions

1. **UUID Primary Keys**: For distributed systems and security
2. **Soft Deletes**: Using status fields rather than deletion
3. **Timestamps**: created_at and updated_at on all tables
4. **Foreign Keys**: Enforced referential integrity
5. **Indexes**: On all foreign keys and frequently queried fields

### Important Constraints

- Company-level data isolation
- Unique SKUs globally
- Unique order numbers
- Email unique per company for customers
- Variant option value combinations

## API Design

### RESTful Principles

- Resource-based URLs
- HTTP verbs (GET, POST, PUT, DELETE)
- JSON request/response
- Consistent error handling

### Authentication (Future)

- Keycloak OIDC tokens
- JWT validation middleware
- Role-based access control
- Company-scoped permissions

### Endpoints Structure

```
/api/v1
  /companies        # Company management
  /customers        # Customer management
  /products         # Product catalog
  /orders           # Order management
  /me               # Customer portal (auth required)
```

## Extensibility Points

### 1. Stock Connectors

Add new connectors by implementing the `Connector` interface:

```go
type WMSConnector struct {
    client *wms.Client
}

func (c *WMSConnector) GetAvailableStock(...) {...}
// Implement other methods
```

### 2. Payment Providers

Future payment module will follow similar pattern:

```go
type PaymentProvider interface {
    CreatePayment(...) (*Payment, error)
    CapturePayment(...) error
    RefundPayment(...) error
}
```

### 3. Shipping Providers

Future shipping integration:

```go
type ShippingProvider interface {
    GetRates(...) ([]Rate, error)
    CreateShipment(...) (*Shipment, error)
    TrackShipment(...) (*TrackingInfo, error)
}
```

## Security Considerations

### Current Implementation

1. **SQL Injection**: Protected via parameterized queries (SQLC)
2. **Input Validation**: HTTP handler level validation
3. **Error Handling**: No sensitive data in error messages
4. **Dependencies**: Regular updates required

### Future Implementation

1. **Authentication**: Keycloak OIDC integration
2. **Authorization**: Role-based access control
3. **Rate Limiting**: API rate limits per user/IP
4. **Audit Logging**: Track all data modifications
5. **Encryption**: TLS in production, encrypted sensitive fields

## Deployment

### Development

```bash
make dev           # Start infrastructure
make migrate-up    # Apply migrations
make run           # Start API server
```

### Production Considerations

1. **Database**: Managed PostgreSQL (AWS RDS, GCP Cloud SQL)
2. **Cache**: Managed Redis (AWS ElastiCache, GCP Memorystore)
3. **Scaling**: Horizontal scaling of API servers
4. **Monitoring**: Prometheus metrics, distributed tracing
5. **Logging**: Centralized log aggregation
6. **Backups**: Automated database backups
7. **Secrets**: Vault or cloud secret management

## Performance Optimization

### Database

- Proper indexes on foreign keys and query fields
- Connection pooling (configured in database package)
- Query optimization via EXPLAIN ANALYZE
- Read replicas for heavy read operations

### Caching Strategy

- Redis for session data
- Product catalog caching
- Price calculation results (with TTL)
- Stock level caching (short TTL)

### API

- Request timeout middleware
- Response compression
- Pagination for list endpoints
- Field filtering (future)

## Testing Strategy

### Unit Tests

- Service layer business logic
- Pricing engine calculations
- Stock connector implementations

### Integration Tests

- Database operations
- HTTP handlers with test DB
- End-to-end flows

### Load Tests

- Concurrent order creation
- Product listing performance
- Price calculation under load

## Future Enhancements

1. **Authentication**: Complete Keycloak integration
2. **Admin Panel**: Management interface
3. **Advanced Pricing**: Time-based rules, customer-specific pricing
4. **Inventory**: Real-time stock updates, low stock alerts
5. **Fulfillment**: Shipping integration, tracking
6. **Invoicing**: Invoice generation, payment tracking
7. **Reporting**: Sales reports, analytics
8. **Multi-language**: i18n support
9. **Multi-currency**: Real-time exchange rates
10. **GraphQL API**: Alternative to REST

## Maintenance

### Regular Tasks

- Dependency updates
- Security patches
- Database backups verification
- Performance monitoring
- Log analysis

### Database Migrations

- Always create migrations with rollback
- Test migrations on staging first
- Zero-downtime deployments when possible
- Keep migrations in version control
