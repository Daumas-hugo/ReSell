package pricing

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// CustomerType represents the type of customer
type CustomerType string

const (
	CustomerTypeB2C CustomerType = "B2C"
	CustomerTypeB2B CustomerType = "B2B"
)

// SalesChannel represents a sales channel
type SalesChannel string

const (
	SalesChannelWeb         SalesChannel = "web"
	SalesChannelPortal      SalesChannel = "portal"
	SalesChannelMobile      SalesChannel = "mobile"
	SalesChannelMarketplace SalesChannel = "marketplace"
	SalesChannelAgency      SalesChannel = "agency"
)

// PriceRequest contains all parameters for price calculation
type PriceRequest struct {
	VariantID     uuid.UUID
	CompanyID     uuid.UUID
	CustomerType  CustomerType
	SalesChannel  SalesChannel
	Quantity      int
	CouponCodes   []string
	Country       string
	CustomerID    *uuid.UUID
}

// PriceResult contains the calculated price with details
type PriceResult struct {
	BasePriceHT     decimal.Decimal
	UnitPriceHT     decimal.Decimal
	UnitPriceTTC    decimal.Decimal
	TotalHT         decimal.Decimal
	TotalTTC        decimal.Decimal
	TaxRate         decimal.Decimal
	TaxAmount       decimal.Decimal
	DiscountsApplied []DiscountDetail
	TaxDetails      TaxDetail
}

// DiscountDetail represents a single discount applied
type DiscountDetail struct {
	Type        string
	Name        string
	Description string
	Amount      decimal.Decimal
	Percentage  decimal.Decimal
}

// TaxDetail represents tax calculation details
type TaxDetail struct {
	Rate    decimal.Decimal
	Amount  decimal.Decimal
	Country string
}

// Engine is the pricing calculation engine
type Engine struct {
	// Dependencies would be injected here
}

// NewEngine creates a new pricing engine
func NewEngine() *Engine {
	return &Engine{}
}

// CalculatePrice calculates the final price based on all rules
func (e *Engine) CalculatePrice(ctx context.Context, req PriceRequest) (*PriceResult, error) {
	// This is a simplified implementation
	// In production, this would:
	// 1. Get base price from variant
	// 2. Apply sales channel rules
	// 3. Apply pricing rules (B2B discounts, volume discounts)
	// 4. Apply promotions
	// 5. Apply coupons
	// 6. Calculate taxes
	
	result := &PriceResult{
		BasePriceHT:      decimal.NewFromFloat(100.00), // Mock base price
		DiscountsApplied: []DiscountDetail{},
	}

	// Start with base price
	currentPrice := result.BasePriceHT

	// Apply quantity
	totalHT := currentPrice.Mul(decimal.NewFromInt(int64(req.Quantity)))

	// Apply discounts (simplified)
	discount := e.calculateDiscounts(ctx, req, currentPrice)
	if !discount.IsZero() {
		totalHT = totalHT.Sub(discount)
		result.DiscountsApplied = append(result.DiscountsApplied, DiscountDetail{
			Type:   "volume",
			Name:   "Volume Discount",
			Amount: discount,
		})
	}

	// Calculate tax
	taxRate := e.getTaxRate(ctx, req.CompanyID, req.Country)
	taxAmount := totalHT.Mul(taxRate).Div(decimal.NewFromInt(100))
	totalTTC := totalHT.Add(taxAmount)

	result.UnitPriceHT = currentPrice
	result.UnitPriceTTC = currentPrice.Add(currentPrice.Mul(taxRate).Div(decimal.NewFromInt(100)))
	result.TotalHT = totalHT
	result.TotalTTC = totalTTC
	result.TaxRate = taxRate
	result.TaxAmount = taxAmount
	result.TaxDetails = TaxDetail{
		Rate:    taxRate,
		Amount:  taxAmount,
		Country: req.Country,
	}

	return result, nil
}

// calculateDiscounts calculates applicable discounts
func (e *Engine) calculateDiscounts(ctx context.Context, req PriceRequest, basePrice decimal.Decimal) decimal.Decimal {
	// Simplified discount calculation
	// In production, this would query pricing_rules and promotions tables
	
	discount := decimal.Zero
	
	// Example: 10% discount for B2B customers
	if req.CustomerType == CustomerTypeB2B {
		discount = basePrice.Mul(decimal.NewFromFloat(0.10))
	}
	
	// Example: Volume discount
	if req.Quantity >= 10 {
		volumeDiscount := basePrice.Mul(decimal.NewFromFloat(0.05))
		if volumeDiscount.GreaterThan(discount) {
			discount = volumeDiscount
		}
	}
	
	return discount.Mul(decimal.NewFromInt(int64(req.Quantity)))
}

// getTaxRate gets the applicable tax rate
func (e *Engine) getTaxRate(ctx context.Context, companyID uuid.UUID, country string) decimal.Decimal {
	// Simplified tax rate lookup
	// In production, this would query the tax_rates table
	
	// Default VAT rate for EU
	if country == "FR" {
		return decimal.NewFromFloat(20.00)
	}
	
	return decimal.NewFromFloat(19.00) // Default
}

// PricingRule represents a pricing rule configuration
type PricingRule struct {
	ID            uuid.UUID
	CompanyID     uuid.UUID
	Name          string
	CustomerType  *CustomerType
	ChannelID     *uuid.UUID
	DiscountType  string
	DiscountValue decimal.Decimal
	MinQuantity   *int
	Priority      int
	ValidFrom     *time.Time
	ValidTo       *time.Time
	Status        string
}

// Promotion represents a promotion configuration
type Promotion struct {
	ID              uuid.UUID
	CompanyID       uuid.UUID
	Name            string
	Code            *string
	Type            string
	DiscountType    string
	DiscountValue   decimal.Decimal
	MinOrderAmount  *decimal.Decimal
	MaxUses         *int
	UsesCount       int
	ValidFrom       time.Time
	ValidTo         time.Time
	Status          string
}
