package r

import (
	"fmt"

	"github.com/cachesdev/souuup/u"
)

// MinN validates if a numeric value is at least n.
//
// Example:
//
//	// Validate that age is at least 18
//	ageField := u.Field(25, u.MinN(18))
func MinN[T u.Numeric](n T) u.NumericRule[T] {
	return func(fd u.FieldState[T]) error {
		if fd.Value < n {
			return fmt.Errorf("value is %v, but needs to be at least %v", fd.Value, n)
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
func MaxN[T u.Numeric](n T) u.NumericRule[T] {
	return func(fd u.FieldState[T]) error {
		if fd.Value > n {
			return fmt.Errorf("value is %v, but needs to be at most %v", fd.Value, n)
		}
		return nil
	}
}
