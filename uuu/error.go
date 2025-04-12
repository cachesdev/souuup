package u

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/k0kubun/pp/v3"
)

// RuleError represents a single validation rule failure.
type RuleError = string

// RuleError represents a slice of rule validation failures.
type RuleErrors = []RuleError

// FieldErrors represents a collection of validation rule failures for a field.
type FieldErrors = map[string]RuleErrors

// ValidationError represents the tree of nested validation errors in a field.
type ValidationError struct {
	Errors       FieldErrors
	NestedErrors map[string]*ValidationError
}

// MarshalJSON implements the json.Marshaler interface for ValidationError.
// It creates a flattened JSON representation where field names are directly mapped
// to their nested structure, with direct errors stored in an "errors" field.
func (ve *ValidationError) MarshalJSON() ([]byte, error) {
	errorMap := ve.ToMap()
	if errorMap == nil {
		return []byte("{}"), nil
	}

	pp.Println(errorMap)

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
