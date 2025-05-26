# Souuup

A robust, type-safe validation library for Go. In any good soup, you need the right balance of ingredients; not too much or too little. Souuup helps you ensure your data has just the right shape and properties, so it too can be tasty.

## Overview

Souuup is a flexible validation framework that lets you easily validate complex data structures with a clear, composable API. It provides type-safe validation using generics with detailed error reporting.

## Features

- 🔍 **Type-safe validation** - Uses generics for compile-time type checking
- 🧩 **Composable rules** - Mix and match validation rules for your specific needs
- 🌳 **Nested validation** - Validate complex, nested data structures
- 🚦 **Detailed error reporting** - Get comprehensive error messages with exact locations

## Quick Start

```go
package main

import (
    "fmt"

    u "github.com/cachesdev/souuup/uuu"
)

// These would be types from your HTTP handler, or database, etc
type User struct {
	Name string
	Age  int
}

type Address struct {
	City    string
	Country string
}

func main() {
	user := User{Name: "John Smith", Age: 27}
	addr := Address{City: "London", Country: "UK"}

	// Create a validator with the schema. You can also extract the schema definition, if you prefer.
	validator := u.NewSouuup(u.Schema{
		"username": u.Field(user.Name, u.MinS(3), u.MaxS(20)),
		"age":      u.Field(user.Age, u.MinN(18), u.MaxN(120)),
		"address": u.Schema{
			"city":    u.Field(addr.City, u.MinS(2)),
			"country": u.Field(addr.Country, u.NotZero),
		},
	})

	// Validate the data
	err := validator.Validate()
	if err != nil {
		fmt.Println("Validation failed:", err)
		return
	}

	fmt.Println("Validation succeeded!")
}
```

## Error Handling

Souuup provides detailed error information, making it easy to identify exactly which fields failed validation and why:

```go
// Example of validation error output (JSON)
{
    "username": {
        "errors": ["length is 2, but needs to be at least 3"]
    },
    "address": {
        "city": {
            "errors": ["length is 1, but needs to be at least 2"]
        }
    }
}
```

## Creating Custom Rules

You can easily create custom validation rules:

```go
// Create a custom rule for email validation
emailRule := func(fd u.FieldState[string]) error {
    if !strings.Contains(fd.value, "@") {
        return fmt.Errorf("must be a valid email address")
    }
    return nil
}

// Use the custom rule
u.Field("user@example.com", emailRule)
```

### Nested Schemas

```go
userSchema := u.Schema{
    "profile": u.Schema{
        "personal": u.Schema{
            "firstName": u.Field("John", u.MinS(2)),
            "lastName":  u.Field("Doe", u.MinS(2)),
        },
        "contact": u.Schema{
            "email": u.Field("john.doe@example.com"),
            "phone": u.Field("+44123456789"),
        },
    },
}
```

### Reusable Validators

If you would like to create validators that you can reuse across multiple handlers, you can create a wrapper function and provide it to your handlers

```go
func ValidateUser(user User) *u.Souuup {
	validator := u.NewSouuup(u.Schema{
		"username": u.Field(user.Name, u.MinS(3), u.MaxS(20)),
		"age":      u.Field(user.Age, u.MinN(18), u.MaxN(120)),
		"address": u.Schema{
			"city":    u.Field(addr.City, u.MinS(2)),
			"country": u.Field(addr.Country, u.NotZero),
		},
	})

    return validator
```

## License

This project is licensed under the terms of the LICENSE file included in the repository.
