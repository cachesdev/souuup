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
//	ageField := u.Field(25, r.MinN(18))
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
//	ageField := u.Field(25, r.MaxN(120))
func MaxN[T u.Numeric](n T) u.NumericRule[T] {
	return func(fd u.FieldState[T]) error {
		if fd.Value > n {
			return fmt.Errorf("value is %v, but needs to be at most %v", fd.Value, n)
		}
		return nil
	}
}

// Gt validates if a numeric value is greater than n.
//
// Example:
//
//	// Validate that age is greater than 18
//	ageField := u.Field(25, r.Gt(18))
func Gt[T u.Numeric](n T) u.NumericRule[T] {
	return func(fd u.FieldState[T]) error {
		if fd.Value <= n {
			return fmt.Errorf("value is %v, but needs to be greater than %v", fd.Value, n)
		}
		return nil
	}
}

// Gte validates if a numeric value greater than or equal to n. It is an alias for MinN.
//
// Example:
//
//	// Validate that age is greater than or equal to 18
//	ageField := u.Field(25, r.Gte(18))
func Gte[T u.Numeric](n T) u.NumericRule[T] {
	return MinN(n)
}

// Lt validates if a numeric value is less than n.
//
// Example:
//
//	// Validate that age is less than 120
//	ageField := u.Field(25, r.Lt(120))
func Lt[T u.Numeric](n T) u.NumericRule[T] {
	return func(fd u.FieldState[T]) error {
		if fd.Value >= n {
			return fmt.Errorf("value is %v, but needs to be less than %v", fd.Value, n)
		}
		return nil
	}
}

// Lte validates if a numeric value less than or equal to n. It is an alias for MaxN.
//
// Example:
//
//	// Validate that age is less than or equal to 120
//	ageField := u.Field(25, r.Lte(120))
func Lte[T u.Numeric](n T) u.NumericRule[T] {
	return MaxN(n)
}

// NeqN validates if a numeric value is not equal to n.
//
// Example:
//
//	// Validate that multi-buy quantity is not 1
//	cartSizeField := u.Field(5, r.NeqN(1))
func NeqN[T u.Numeric](n T) u.NumericRule[T] {
	return func(fd u.FieldState[T]) error {
		if fd.Value == n {
			return fmt.Errorf("value is %v, but needs to not equal to %v", fd.Value, n)
		}
		return nil
	}
}
