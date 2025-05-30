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
//	emailField := u.Field("user@example.com", u.NotZero)
//
//	// Works with any comparable type
//	ageField := u.Field(25, r.NotZero)
func NotZero[T comparable](fd u.FieldState[T]) error {
	var zero T
	if fd.Value == zero {
		return fmt.Errorf("value is required but has zero value")
	}
	return nil
}
