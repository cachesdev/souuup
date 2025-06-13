// Package main demonstrates complex validation scenarios with Souuup
package main

import (
	"fmt"
	"time"

	"github.com/cachesdev/souuup/r"
	"github.com/cachesdev/souuup/u"
)

// Complex data structures for validation
type Order struct {
	OrderID      string
	CustomerID   string
	OrderDate    time.Time
	ShipDate     *time.Time
	Items        []OrderItem
	ShippingInfo Address
	PaymentInfo  PaymentInfo
	TotalAmount  float64
	Status       string
}

type OrderItem struct {
	ProductID   string
	Quantity    int
	UnitPrice   float64
	Discount    float64
	Description string
}

type Address struct {
	Street        string
	City          string
	State         string
	Country       string
	PostalCode    string
	IsResidential bool
}

type PaymentInfo struct {
	Method          string
	CardLastFour    string
	ExpirationMonth int
	ExpirationYear  int
	Paid            bool
}

// Custom validation rules
func FutureDate(fs u.FieldState[time.Time]) error {
	if fs.Value.Before(time.Now()) {
		return fmt.Errorf("date must be in the future")
	}
	return nil
}

func PastDate(fs u.FieldState[time.Time]) error {
	if fs.Value.After(time.Now()) {
		return fmt.Errorf("date must be in the past")
	}
	return nil
}

func ValidCardExpiration(month, year int) bool {
	now := time.Now()
	currentYear, currentMonth := now.Year(), int(now.Month())

	// If expiration year is in the past, it's invalid
	if year < currentYear {
		return false
	}

	// If expiration year is current year but month is in the past, it's invalid
	if year == currentYear && month < currentMonth {
		return false
	}

	return true
}

func ValidPaymentMethod(fs u.FieldState[string]) error {
	validMethods := map[string]bool{
		"credit_card":   true,
		"debit_card":    true,
		"paypal":        true,
		"bank_transfer": true,
	}

	if !validMethods[fs.Value] {
		return fmt.Errorf("invalid payment method: must be one of credit_card, debit_card, paypal, or bank_transfer")
	}

	return nil
}

func HasMinimumValue(minValue float64) u.Rule[OrderItem] {
	return func(fs u.FieldState[OrderItem]) error {
		item := fs.Value
		totalValue := item.UnitPrice*float64(item.Quantity) - item.Discount

		if totalValue < minValue {
			return fmt.Errorf("total item value (%.2f) is below minimum threshold of %.2f", totalValue, minValue)
		}

		return nil
	}
}

func main() {
	fmt.Println("Complex Validation Example")
	fmt.Println("=========================")

	// Sample data
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	order := Order{
		OrderID:    "ORD-12345",
		CustomerID: "CUST-789",
		OrderDate:  now,
		ShipDate:   &tomorrow,
		Items: []OrderItem{
			{
				ProductID:   "PROD-001",
				Quantity:    2,
				UnitPrice:   49.99,
				Discount:    5.00,
				Description: "Wireless Headphones",
			},
			{
				ProductID:   "PROD-002",
				Quantity:    1,
				UnitPrice:   999.99,
				Discount:    0.00,
				Description: "Smartphone",
			},
		},
		ShippingInfo: Address{
			Street:        "123 Main St",
			City:          "Boston",
			State:         "MA",
			Country:       "USA",
			PostalCode:    "02108",
			IsResidential: true,
		},
		PaymentInfo: PaymentInfo{
			Method:          "credit_card",
			CardLastFour:    "1234",
			ExpirationMonth: 12,
			ExpirationYear:  time.Now().Year() + 2,
			Paid:            true,
		},
		TotalAmount: 1094.97,
		Status:      "processing",
	}

	// Create complex nested validation schema
	orderSchema := u.Schema{
		"orderID":    u.Field(order.OrderID, r.NotZero, r.MinS(5)),
		"customerID": u.Field(order.CustomerID, r.NotZero),
		"orderDate":  u.Field(order.OrderDate, PastDate),
		"shipDate":   u.Field(*order.ShipDate, FutureDate),
		"items": u.Field(order.Items,
			r.MinLen[OrderItem](1),         // At least one item required
			r.MaxLen[OrderItem](10),        // Maximum 10 items allowed
			r.Some(HasMinimumValue(100.0)), // At least one item must be worth $100 or more,
		),
		"shippingInfo": u.Schema{
			"street":     u.Field(order.ShippingInfo.Street, r.NotZero, r.MinS(5)),
			"city":       u.Field(order.ShippingInfo.City, r.NotZero, r.MinS(2)),
			"state":      u.Field(order.ShippingInfo.State, r.NotZero, r.MinS(2)),
			"country":    u.Field(order.ShippingInfo.Country, r.NotZero, r.MinS(2)),
			"postalCode": u.Field(order.ShippingInfo.PostalCode, r.NotZero),
		},
		"paymentInfo": u.Schema{
			"method": u.Field(order.PaymentInfo.Method, ValidPaymentMethod),
			"cardLastFour": u.Field(order.PaymentInfo.CardLastFour, r.NotZero, func(fs u.FieldState[string]) error {
				if len(fs.Value) != 4 {
					return fmt.Errorf("card last four must be exactly 4 digits")
				}
				return nil
			}),
			"expiration": u.Field(true, func(fs u.FieldState[bool]) error {
				if !ValidCardExpiration(order.PaymentInfo.ExpirationMonth, order.PaymentInfo.ExpirationYear) {
					return fmt.Errorf("card expiration date is invalid")
				}
				return nil
			}),
		},
		"totalAmount": u.Field(order.TotalAmount, r.MinN(0.0)),
		"status": u.Field(order.Status, func(fs u.FieldState[string]) error {
			validStatuses := map[string]bool{
				"pending":    true,
				"processing": true,
				"shipped":    true,
				"delivered":  true,
				"cancelled":  true,
			}

			if !validStatuses[fs.Value] {
				return fmt.Errorf("invalid status")
			}
			return nil
		}),
	}

	// Create validator
	s := u.NewSouuup(orderSchema)

	// Validate order
	err := s.Validate()
	if err != nil {
		fmt.Printf("Order validation failed: %s\n", err)
		return
	}

	fmt.Println("✅ Order validated successfully!")

	// Create an invalid order to demonstrate validation failures
	invalidShipDate := now.Add(-24 * time.Hour) // yesterday (invalid ship date)
	invalidOrder := order
	invalidOrder.ShipDate = &invalidShipDate
	invalidOrder.Items = []OrderItem{
		{
			ProductID:   "PROD-001",
			Quantity:    2,
			UnitPrice:   49.99,
			Discount:    5.00,
			Description: "Wireless Headphones",
		},
		{
			ProductID:   "PROD-002",
			Quantity:    1,
			UnitPrice:   9.99, // Too cheap to satisfy the minimum value rule
			Discount:    0.00,
			Description: "Budget Earbuds",
		},
	}
	invalidOrder.PaymentInfo.Method = "bitcoin" // Invalid payment method

	// Create validation schema for invalid order
	invalidOrderSchema := u.Schema{
		"orderID":    u.Field(invalidOrder.OrderID, r.NotZero, r.MinS(5)),
		"customerID": u.Field(invalidOrder.CustomerID, r.NotZero),
		"orderDate":  u.Field(invalidOrder.OrderDate, PastDate),
		"shipDate":   u.Field(*invalidOrder.ShipDate, FutureDate),
		"items": u.Schema{
			"count": u.Field(len(invalidOrder.Items), r.MinN(1)),
			"item0": u.Schema{
				"productID": u.Field(invalidOrder.Items[0].ProductID, r.NotZero),
				"quantity":  u.Field(invalidOrder.Items[0].Quantity, r.MinN(1)),
				"unitPrice": u.Field(invalidOrder.Items[0].UnitPrice, r.MinN(0.01)),
			},
		},
		"paymentInfo": u.Schema{
			"method": u.Field(invalidOrder.PaymentInfo.Method, ValidPaymentMethod),
		},
	}

	// Create validator for invalid order
	invalid := u.NewSouuup(invalidOrderSchema)

	// Validate invalid order
	fmt.Println("\nValidating an invalid order...")
	invalidErr := invalid.Validate()
	if invalidErr != nil {
		fmt.Printf("❌ Invalid order validation failed as expected: %s\n", invalidErr)
		return
	}

	fmt.Println("⚠️ Invalid order unexpectedly passed validation!")
}
