package u

import (
	"encoding/json"
	"fmt"
)

// RuleError represents a single validation rule failure.
type RuleError string

// Needed for internal usage in Souuup.Validate()
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

// Alias for brevity when returning from ValidationError.ToMap()
type ToMapResult = map[string]map[string]any

// Helper to avoid nil maps
func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors:       make(FieldsErrorMap),
		NestedErrors: make(map[string]*ValidationError),
	}
}

// Adds an error to the ValidationError's top level Errors, stored against the field
func (ve *ValidationError) AddError(field string, err error) {
	ve.Errors[field] = append(ve.Errors[field], RuleError(err.Error()))
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

// ToMap converts a ValidationError to a map representation.
// This provides a more efficient way to create the flattened error structure.
// It recursively processes the entire validation error tree and returns a map where
// field names are directly mapped to objects containing:
// - "errors": array of direct errors for the field
// - Other keys: nested validation structures
//
// INFO: Does many recursive calls. maybe performance issues?
func (ve *ValidationError) ToMap() ToMapResult {
	if !ve.HasErrors() {
		return nil
	}

	result := make(ToMapResult, len(ve.Errors)+len(ve.NestedErrors))

	// Add direct field errors
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
				// Manually copy each key-value pair because maps.Copy isn't type smart enough
				for k, v := range nestedMap {
					existing[k] = v
				}
			} else {
				// No existing entry, add the nested map directly
				// Copy all nested values since the type inference isn't smart enough to directly assign nestedMap
				result[nestedField] = make(map[string]any)
				for k, v := range nestedMap {
					result[nestedField][k] = v
				}
			}
		}
	}

	return result
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

// GetOrCreateNested returns a nested ValidationError for a field, creating it if necessary.
func (ve *ValidationError) GetOrCreateNested(field string) *ValidationError {
	if _, exists := ve.NestedErrors[field]; !exists {
		ve.NestedErrors[field] = NewValidationError()
	}
	return ve.NestedErrors[field]
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
