package u

// RuleError represents a single validation rule failure.
type RuleError string

func (fe RuleError) Error() string {
	return string(fe)
}

// FieldErrors represents a collection of validation rule failures for a field.
type FieldErrors = map[string][]RuleError

// ValidationError represents the tree of nested validation errors in a field.
type ValidationError struct {
	Errors       FieldErrors
	NestedErrors map[string]*ValidationError
}

func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors:       make(FieldErrors),
		NestedErrors: make(map[string]*ValidationError),
	}
}

func (ve *ValidationError) AddError(field string, err RuleError) {
	ve.Errors[field] = append(ve.Errors[field], err)
}

// GetOrCreateNested returns a nested ValidationError for a field, creating it if necessary.
func (ve *ValidationError) GetOrCreateNested(field string) *ValidationError {
	if _, exists := ve.NestedErrors[field]; !exists {
		ve.NestedErrors[field] = NewValidationError()
	}
	return ve.NestedErrors[field]
}

// HasErrors returns true if there are any errors at any level.
func (ve *ValidationError) HasErrors() bool {
	if len(ve.Errors) > 0 {
		return true
	}

	for _, nestedErr := range ve.NestedErrors {
		if nestedErr.HasErrors() {
			return true
		}
	}

	return false
}

// TODO: Currently no-op
func (ve *ValidationError) Error() string {
	if !ve.HasErrors() {
		return ""
	}

	// beep boop implementation details
	return "validation failed"
}
