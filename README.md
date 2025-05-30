# Souuup

A robust, type-safe validation library for Go. In any good soup, you need the right ingredients. Souuup helps you ensure your data has just the right shape and properties, in a type safe way.

## Overview

Souuup is a flexible validation framework that lets you easily validate complex data structures with a clear, composable API. It provides type-safe validation using generics with detailed error reporting.

## Features

- 🔍 **Type-safe validation** - Uses generics for compile-time type checking
- 🧩 **Composable rules** - Mix and match validation rules for your specific needs
- 🌳 **Nested validation** - Validate complex, nested data structures
- 🚦 **Detailed error reporting** - Get comprehensive error messages with the same shape as your schema

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/cachesdev/souuup/u"
    "github.com/cachesdev/souuup/r"
)

type User struct {
	Name string
	Age  int
	Interests []string
}

type Address struct {
	City    string
	Country string
}

func main() {
	user := User{
		Name: "John Smith",
		Age: 27,
		Interests: []string{"reading", "cycling", "photography"},
	}
	addr := Address{City: "London", Country: "UK"}

	// Create a validator with the schema or extract the schema definition, if you prefer.
	s := u.NewSouuup(u.Schema{
		"username": u.Field(user.Name, r.MinS(3), r.MaxS(20)),
		"age":      u.Field(user.Age, r.MinN(18), r.MaxN(120)),
				"interests": u.Field(user.Interests,
			u.MinLen[string](1),    // At least one interest required
			u.Every(u.MinS(3)),        // Each interest must be at least 3 characters
		),
		"address": u.Schema{
			"city":    u.Field(addr.City, r.MinS(2)),
			"country": u.Field(addr.Country, r.NotZero),
		},
	})

	// Validate the data
	err := s.Validate()
	if err != nil {
		fmt.Println("Validation failed:", err)
		return
	}

	fmt.Println("Validation succeeded!")
}
```

## Error Handling

Souuup provides detailed error information, making it easy to identify exactly which fields failed validation and why:

```json
// Example of validation error output (JSON)
{
  "username": {
    "errors": ["length is 2, but needs to be at least 3"]
  },
  "interests": {
    "errors": ["2 elements failed validation
      [0]: length is 1, but needs to be at least 3
      [1]: length is 2, but needs to be at least 3"]
    }
  }
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
    if !strings.Contains(fd.Value, "@") {
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
            "firstName": u.Field("John", r.MinS(2)),
            "lastName":  u.Field("Doe", r.MinS(2)),
        },
        "contact": u.Schema{
            "email": u.Field("john.doe@example.com"),
            "phone": u.Field("+44123456789"),
        },
    },
}
```

### Reusable Validators

If you would like to create validators that you can reuse, you can create wrappers or methods and provide it where needed:

```go
type Address struct {
	City string
	Country string
}

type User struct {
	Name string
	Age int
	Interests []string
	Address Address
}
func (user User) Validate() error {
	s := u.NewSouuup(u.Schema{
		"username": u.Field(user.Name, r.MinS(3), r.MaxS(20)),
		"age":      u.Field(user.Age, r.MinN(18), r.MaxN(120)),
		"interests": u.Field(user.Interests,
            u.MinLen[string](1),
            u.Every(u.MinS(3)),
        ),
		"address": u.Schema{
			"city":    u.Field(user.Address.City, r.MinS(2)),
			"country": u.Field(user.Address.Country, r.NotZero),
		},
	})

    return s.Validate()
}

func main() {
	user := User{
		Name: "Juan",
		Age: 25,
		Address: Address{
			City: "London",
		}
	}

	err := user.Validate()
	fmt.Println(err)
	//{
	//  "address": {
	//    "country": {
	//        "errors": ["field country cannot be empty"]
	//    }
	//  }
	//}
```

## License

This project is licensed under the terms of the LICENSE file included in the repository.
