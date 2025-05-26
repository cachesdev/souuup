package u

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// Numeric is a constraint that permits any numeric type: integer or float.
type Numeric interface {
	constraints.Float | constraints.Integer
}

// Rule is a generic validation rule that takes a field state and returns an error if validation fails.
// It is the fundamental building block of the validation system.
type Rule[T any] = func(FieldState[T]) error

// StringRule is a specialised rule type for string validation.
type StringRule = Rule[string]

// NumericRule is a specialised rule type for numeric validation.
type NumericRule[T Numeric] = Rule[T]

// MinS validates if a string's length is at least n characters.
//
// Example:
//
//	// Validate that a name is at least 2 characters long
//	nameField := u.Field("John", u.MinS(2))
func MinS(n int) StringRule {
	return func(fd FieldState[string]) error {
		if len(fd.value) < n {
			return fmt.Errorf("length is %d, but needs to be at least %d", len(fd.value), n)
		}
		return nil
	}
}

// MinN validates if a numeric value is at least n.
//
// Example:
//
//	// Validate that age is at least 18
//	ageField := u.Field(25, u.MinN(18))
func MinN[T Numeric](n T) NumericRule[T] {
	return func(fd FieldState[T]) error {
		if fd.value < n {
			return fmt.Errorf("value is %v, but needs to be at least %v", fd.value, n)
		}
		return nil
	}
}

// MaxS validates if a string's length is at most n characters.
//
// Example:
//
//	// Validate that a username is at most 20 characters long
//	usernameField := u.Field("john doe", u.MaxS(20))
func MaxS(n int) StringRule {
	return func(fd FieldState[string]) error {
		if len(fd.value) > n {
			return fmt.Errorf("length is %d, but needs to be at most %d", len(fd.value), n)
		}
		return nil
	}
}

// MaxN validates if a numeric value is at most n.
//
// Example:
//
//	// Validate that age is at most 120
//	ageField := u.Field(25, u.MaxN(120))
func MaxN[T Numeric](n T) NumericRule[T] {
	return func(fd FieldState[T]) error {
		if fd.value > n {
			return fmt.Errorf("value is %v, but needs to be at most %v", fd.value, n)
		}
		return nil
	}
}

// NotZero validates that a value is not the zero value for its type.
// This is useful for required fields.
//
// Example:
//
//	// Validate that email is not empty
//	emailField := u.Field("user@example.com", u.NotZero)
//
//	// Works with any comparable type
//	ageField := u.Field(25, u.NotZero)
func NotZero[T comparable](fd FieldState[T]) error {
	var zero T
	if fd.value == zero {
		return fmt.Errorf("value is required but has zero value")
	}
	return nil
}
