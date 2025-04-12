package u

// Rule is a generic validation rule for any type
type Rule[T any] = func(FieldState[T]) bool

// StringRule is for string-specific validation rules
type StringRule = Rule[string]

// NumericRule is for numeric validation rules
type NumericRule[T Numeric] = Rule[T]

// MinS validates if a string is shorter than n characters
func MinS(n int) StringRule {
	return func(fd FieldState[string]) bool {
		return len(fd.value) < n
	}
}

// MinN validates if a numeric value is less than n
func MinN[T Numeric](n T) NumericRule[T] {
	return func(fd FieldState[T]) bool {
		return fd.value < n
	}
}
