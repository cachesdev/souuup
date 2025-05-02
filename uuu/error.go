package u

import (
	"encoding/json"
	"fmt"
	"maps"
)

// RuleError represents a single validation rule failure.
type RuleError string

// Neded for internal usage in Souup.Validate()
func (re RuleError) Error() string {
	return string(re)
}

// RuleError represents a slice of rule validation failures.
type RuleErrors = []RuleError

// FieldsErrorMap represents many fields K with their V RuleErrors.
type FieldsErrorMap = map[string]RuleErrors

// NestedErrorsMap represents the collection of fields under the current one,
// and each of their validation errors
type NestedErrorsMap = map[string]*ValidationError

// ValidationError represents the tree of nested validation errors in a field.
type ValidationError struct {
	Errors       FieldsErrorMap
	NestedErrors NestedErrorsMap
	Parent       *ValidationError
}

// MarshalJSON implements the json.Marshaler interface for ValidationError.
// It creates a flattened JSON representation where field names are directly mapped
// to their nested structure, with direct errors stored in an "errors" field.
func (ve *ValidationError) MarshalJSON() ([]byte, error) {
	errorMap := ve.ToMap()
	if errorMap == nil {
		return []byte("{}"), nil
	}

	return json.Marshal(errorMap)
}

// ToMap converts a ValidationError to a map representation.
// This provides a more efficient way to create the flattened error structure.
// It recursively processes the entire validation error tree and returns a map where
// field names are directly mapped to objects containing:
// - "errors": array of direct errors for the field
// - Other keys: nested validation structures
//
// INFO: Does many recursive calls. maybe performance issues?
func (ve *ValidationError) ToMap() map[string]any {
	// recursive
	if !ve.HasErrors() {
		return nil
	}

	result := make(map[string]any, len(ve.Errors)+len(ve.NestedErrors))

	// Add direct field errors
	// state 1: {"field1": {"errors": ["a validation error"]}}
	for field, errors := range ve.Errors {
		result[field] = map[string]any{
			"errors": errors,
		}
	}

	// Add nested field errors
	for nestedField, nestedErr := range ve.NestedErrors {
		if nestedErr.HasErrors() {
			// recursive
			nestedMap := nestedErr.ToMap()

			// Check for existing entry.
			if existing, exists := result[nestedField]; exists {
				// If entry exists as a map (from direct errors), merge the nested data
				// the merge is between keys. same keys will be will replaced!
				if existingMap, ok := existing.(map[string]any); ok {
					maps.Copy(existingMap, nestedMap)
				}
			} else {
				// No existing entry, add the nested map directly
				result[nestedField] = nestedMap
			}
		}
	}

	return result
}

// Helper to avoid nil maps
func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors:       make(FieldsErrorMap),
		NestedErrors: make(map[string]*ValidationError),
	}
}

func (ve *ValidationError) AddError(field string, err error) {
	ve.Errors[field] = append(ve.Errors[field], RuleError(err.Error()))
}

// GetOrCreateNested returns a nested ValidationError for a field, creating it if necessary.
func (ve *ValidationError) GetOrCreateNested(field string) *ValidationError {
	if _, exists := ve.NestedErrors[field]; !exists {
		ve.NestedErrors[field] = NewValidationError()
	}
	return ve.NestedErrors[field]
}

// HasErrors returns true if there are any errors at any level, recursively.
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

func (ve *ValidationError) Error() string {
	if !ve.HasErrors() {
		return ""
	}

	bytes, err := json.Marshal(ve)
	if err != nil {
		fmt.Printf("marshalling error %s", err.Error())
	}

	return string(bytes)
}
