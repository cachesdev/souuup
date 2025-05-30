package u

import (
	"golang.org/x/exp/constraints"
)

// Rule is a generic validation rule that takes a field state and returns an error if validation fails.
// It is the fundamental building block of the validation system.
type Rule[T any] = func(FieldState[T]) error

// Specific rule type aliases:

// Numeric is a constraint that permits any numeric type: integer or float.
type Numeric interface {
	constraints.Float | constraints.Integer
}

// StringRule is a specialised rule type for string validation.
type StringRule = Rule[string]

// NumericRule is a specialised rule type for numeric validation.
type NumericRule[T Numeric] = Rule[T]

// SliceRule is a specialised rule type for slice validation.
type SliceRule[T any] = Rule[[]T]
