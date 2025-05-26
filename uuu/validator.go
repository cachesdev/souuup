// Package u provides a robust validation framework for Go.
//
// Souuup is a flexible, type-safe validation library built with generics.
// It allows developers to validate complex data structures with nested fields
// and custom validation rules. The library provides detailed error reporting
// with exact locations of validation failures.
//
// Basic Example:
//
//	schema := u.Schema{
//		"username": u.Field("john doe", u.MinS(3), u.MaxS(20)),
//		"age":      u.Field(25, u.MinN(18)),
//	}
//
//	uuu := u.NewSouuup(schema)
//	err := uuu.Validate()
//	if err != nil {
//		fmt.Println("Validation failed:", err)
//		return
//	}
//
// The package is designed to be intuitive, composable, and extensible, making it
// suitable for a wide range of validation scenarios from simple form validation
// to complex API request validation.
package u

// FieldTag represents the "key" of a field, and will be used to identify a field on
// an error map and schema
type FieldTag = string

// souuupState contains the internal state for a Souuup validator instance.
// This includes any validation errors that have been accumulated.
type souuupState struct {
	errors *ValidationError
}

// Schema is a map of field tags to validatable entities.
// It can contain both simple fields and nested schemas, allowing for
// the validation of complex, hierarchical data structures.
// Schema itself implements the Validable interface.
type Schema map[FieldTag]Validable

var _ Validable = (*Schema)(nil)

// Souuup is the main validator instance.
// It holds a validation schema and internal state.
type Souuup struct {
	state  *souuupState
	schema Schema
}

// NewSouuup creates a new validator instance with the provided schema.
// This is the main entry point for setting up validation.
//
// Example:
//
//	schema := u.Schema{
//		"username": u.Field("johndoe", u.MinS(3)),
//		"age":      u.Field(25, u.MinN(18)),
//	}
//	uuu := u.NewSouuup(schema)
func NewSouuup(schema Schema) *Souuup {
	return &Souuup{
		state:  &souuupState{NewValidationError()},
		schema: schema,
	}
}

// Validate performs validation against the schema and returns an error if validation fails.
// If validation succeeds, it returns nil.
//
// Example:
//
//	err := uuu.Validate()
//	if err != nil {
//		fmt.Println("Validation failed:", err)
//		return
//	}
func (u *Souuup) Validate() error {
	u.schema.Validate(u.state.errors, "")

	if u.state.errors.HasErrors() {
		return u.state.errors
	}
	return nil
}

// Validate implements the Validable interface for Schema.
// It validates all fields and nested schemas within the current schema,
// adding any validation errors to the provided ValidationError object.
func (s Schema) Validate(ve *ValidationError, _ FieldTag) {
	for tag, fieldOrSchema := range s {
		if schema, ok := fieldOrSchema.(Schema); ok {
			newVe := NewValidationError()
			newVe.Parent = ve
			ve.NestedErrors[tag] = newVe
			schema.Validate(newVe, tag)
		} else {
			field := fieldOrSchema
			field.Validate(ve, tag)
		}
	}
}

// Errors returns all validation errors for this schema.
// It creates a new ValidationError, validates the schema against it,
// and returns the resulting errors.
func (s Schema) Errors() *ValidationError {
	errors := NewValidationError()
	s.Validate(errors, "")
	return errors
}

// Validable is the interface that must be implemented by any validatable entity.
// It provides methods to validate the entity and retrieve validation errors.
// Both Schema and FieldDef implement this interface.
type Validable interface {
	// Validate validates the entity against a ValidationError object
	// and associates any errors with the provided field tag.
	Validate(*ValidationError, FieldTag)

	// Errors returns any validation errors associated with this entity.
	Errors() *ValidationError
}
