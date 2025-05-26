// Package main demonstrates basic usage of Souuup for data validation
package main

import (
	"fmt"
	"strings"

	u "github.com/cachesdev/souuup/uuu"
)

// Custom email validation rule
func EmailRule(fs u.FieldState[string]) error {
	// Use the Value() method to access the field value
	email := fs.Value()
	if !strings.Contains(email, "@") {
		return fmt.Errorf("must be a valid email address")
	}
	return nil
}

// Example user data structure (for demonstration purposes)
type User struct {
	Username  string
	Email     string
	Age       int
	IsActive  bool
	Address   Address
	Interests []string
}

type Address struct {
	Street  string
	City    string
	Country string
	ZipCode string
}

func main() {
	fmt.Println("Basic Souuup Validation Example")
	fmt.Println("===============================")

	// Sample user data
	user := User{
		Username: "john",
		Email:    "johndoe@example.com",
		Age:      25,
		IsActive: true,
		Address: Address{
			Street:  "123 Main St",
			City:    "London",
			Country: "UK",
			ZipCode: "W1A 1AA",
		},
		Interests: []string{"reading", "cycling"},
	}

	// Create a validation schema for the user
	schema := u.Schema{
		"username": u.Field(user.Username, u.MinS(3), u.MaxS(20)),
		"email":    u.Field(user.Email, u.NotZero, EmailRule),
		"age":      u.Field(user.Age, u.MinN(18), u.MaxN(120)),
		"isActive": u.Field(user.IsActive),
		"address": u.Schema{
			"street":  u.Field(user.Address.Street, u.NotZero, u.MinS(5)),
			"city":    u.Field(user.Address.City, u.NotZero, u.MinS(2)),
			"country": u.Field(user.Address.Country, u.NotZero, u.MinS(2)),
			"zipCode": u.Field(user.Address.ZipCode, u.NotZero),
		},
		"interests": u.Field(len(user.Interests), u.MinN(1)), // Validate array length
	}

	// Create validator
	validator := u.NewSouuup(schema)

	// Validate data
	err := validator.Validate()
	if err != nil {
		fmt.Printf("Validation failed: %s\n", err)
		return
	}

	fmt.Println("✅ Valid user data validated successfully!")
	fmt.Println()

	// Example with invalid data
	fmt.Println("Now testing with invalid data...")
	invalidUser := User{
		Username: "j",             // Too short
		Email:    "invalid-email", // Missing @ symbol
		Age:      15,              // Too young
		Address: Address{
			Street: "123", // Too short
			City:   "",    // Empty
		},
	}

	// Create a validation schema for the invalid user
	invalidSchema := u.Schema{
		"username": u.Field(invalidUser.Username, u.MinS(3), u.MaxS(20)),
		"email":    u.Field(invalidUser.Email, u.NotZero, EmailRule),
		"age":      u.Field(invalidUser.Age, u.MinN(18), u.MaxN(120)),
		"address": u.Schema{
			"street":  u.Field(invalidUser.Address.Street, u.NotZero, u.MinS(5)),
			"city":    u.Field(invalidUser.Address.City, u.NotZero, u.MinS(2)),
			"country": u.Field(invalidUser.Address.Country, u.NotZero, u.MinS(2)),
		},
	}

	// Create validator
	invalidValidator := u.NewSouuup(invalidSchema)

	// Validate data
	invalidErr := invalidValidator.Validate()
	if invalidErr != nil {
		fmt.Printf("❌ Invalid user validation failed as expected: %s\n", invalidErr)
		return
	}

	fmt.Println("⚠️ Invalid user data incorrectly passed validation!")
}
