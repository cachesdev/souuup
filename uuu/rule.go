package u

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Float | constraints.Integer
}

// Rule is a generic validation rule for any type
type Rule[T any] = func(FieldState[T]) error

// StringRule is for string-specific validation rules
type StringRule = Rule[string]

// NumericRule is for numeric validation rules
type NumericRule[T Numeric] = Rule[T]

// MinS validates if a string is shorter than n characters
func MinS(n int) StringRule {
	return func(fd FieldState[string]) error {
		if len(fd.value) < n {
			return fmt.Errorf("%s is %d long, but needs to be at least %d long", fd.tag, len(fd.value), n)
		}
		return nil
	}
}

// MinN validates if a numeric value is less than n
func MinN[T Numeric](n T) NumericRule[T] {
	return func(fd FieldState[T]) error {
		if fd.value < n {
			return fmt.Errorf("%s is %v, but needs to be at least %v", fd.tag, fd.value, n)
		}
		return nil
	}
}

// MaxS validates if a string is no longer than n characters
func MaxS(n int) StringRule {
	return func(fd FieldState[string]) error {
		if len(fd.value) > n {
			return fmt.Errorf("%s is %d long, but needs to be at most %d long", fd.tag, len(fd.value), n)
		}
		return nil
	}
}

// MaxN validates if a numeric value is not greater than n
func MaxN[T Numeric](n T) NumericRule[T] {
	return func(fd FieldState[T]) error {
		if fd.value > n {
			return fmt.Errorf("%s is %v, but needs to be at most %v", fd.tag, fd.value, n)
		}
		return nil
	}
}

// NotZero validates that a value is not the zero value for its type
func NotZero[T comparable](fd FieldState[T]) error {
	var zero T
	if fd.value == zero {
		return fmt.Errorf("%s is required but has zero value", fd.tag)
	}
	return nil
}
