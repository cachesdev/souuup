package u

import (
	"encoding/json"
	"fmt"
)

// RuleError represents a single validation rule failure.
// It implements the error interface and contains the error message.
type RuleError string

// Error returns the error message for a rule validation failure.
// This implementation satisfies the error interface.
func (re RuleError) Error() string {
	return string(re)
}

// RuleErrors represents a slice of rule validation failures for a single field.
type RuleErrors = []RuleError

// FieldsErrorMap maps field tags to their validation errors.
// This is used to track which fields have validation errors and what those errors are.
type FieldsErrorMap = map[FieldTag]RuleErrors

// NestedErrorsMap maps field tags to nested validation error objects.
// This is used to represent hierarchical validation errors in nested structures.
type NestedErrorsMap = map[FieldTag]*ValidationError

// ValidationError represents the complete tree of validation errors.
// It tracks direct field errors and nested validation errors, forming a tree structure
// that matches the structure of the validated data.
type ValidationError struct {
	// Errors contains direct validation errors for fields at the current level
	Errors FieldsErrorMap

	// NestedErrors contains validation errors for nested fields/structures
	NestedErrors NestedErrorsMap

	// Parent points to the parent ValidationError in the tree, if any
	Parent *ValidationError
}

// ToMapResult is the type returned by ValidationError.ToMap().
// It provides a serialisable representation of validation errors.
type ToMapResult = map[FieldTag]map[string]any

// NewValidationError creates a new ValidationError with initialised maps.
// This helper prevents nil map errors when adding validation failures.
func NewValidationError() *ValidationError {
	return &ValidationError{
		Errors:       make(FieldsErrorMap),
		NestedErrors: make(NestedErrorsMap),
	}
}

// AddError adds a validation error for a specific field tag.
// The error is converted to a RuleError and appended to any existing errors for that field.
func (ve *ValidationError) AddError(tag FieldTag, err error) {
	ve.Errors[tag] = append(ve.Errors[tag], RuleError(err.Error()))
}

// HasErrors returns true if there are any validation errors at any level in the tree.
// It recursively checks nested errors to determine if validation has failed anywhere.
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

// ToMap converts a ValidationError to a map representation suitable for serialisation.
// It recursively processes the entire validation error tree and returns a flattened structure
// where field names are mapped to objects containing:
// - "errors": array of direct errors for the field
// - Other keys: nested validation structures
//
// Example output structure:
//
//	{
//	  "username": {
//	    "errors": ["length is 2, but needs to be at least 3"]
//	  },
//	  "address": {
//	    "city": {
//	      "errors": ["cannot be empty"]
//	    }
//	  }
//	}
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
				result[nestedField] = make(map[FieldTag]any)
				for k, v := range nestedMap {
					result[nestedField][k] = v
				}
			}
		}
	}

	return result
}

// MarshalJSON implements the json.Marshaler interface for ValidationError.
// It creates a JSON representation of the validation errors using the ToMap method.
func (ve *ValidationError) MarshalJSON() ([]byte, error) {
	errorMap := ve.ToMap()
	if errorMap == nil {
		return []byte("null"), nil
	}

	return json.Marshal(errorMap)
}

// GetOrCreateNested returns a nested ValidationError for a field, creating it if necessary.
// This is used when building up validation errors for nested structures.
func (ve *ValidationError) GetOrCreateNested(tag FieldTag) *ValidationError {
	if _, exists := ve.NestedErrors[tag]; !exists {
		ve.NestedErrors[tag] = NewValidationError()
	}
	return ve.NestedErrors[tag]
}

// Error returns a JSON string representation of the validation errors.
// This implementation satisfies the error interface.
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
