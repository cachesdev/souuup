package r

import (
	"fmt"

	"github.com/cachesdev/souuup/u"
)

// NotZero validates that a value is not the zero value for its type.
// This is useful for required fields.
//
// Example:
//
//	// Validate that email is not empty
//	emailField := u.Field("user@example.com", r.NotZero)
//
//	// Works with any comparable type
//	ageField := u.Field(25, r.NotZero)
func NotZero[T comparable](fs u.FieldState[T]) error {
	var zero T
	if fs.Value == zero {
		return fmt.Errorf("value is required but has zero value")
	}
	return nil
}

// SameAs validates that a value is the same as another value of its type.
// This is useful for password/email confirmations.
//
// Example:
//
//	// Validate that passwords match
//	passwordField := u.Field(password, r.SameAs(confirmPassword))
func SameAs[T comparable](other T) u.Rule[T] {
	return func(fs u.FieldState[T]) error {
		if fs.Value != other {
			return fmt.Errorf("%v does not match %v", fs.Value, other)
		}
		return nil
	}
}
