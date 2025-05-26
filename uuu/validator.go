// main package of Soup
package u

// FieldTag represents the "key" of a field, and will be used to identify a field on
// an error map and schema
type FieldTag = string

// Internal State
type souuupState struct {
	errors *ValidationError
}

// Schema of all validable fields/nested fields
// Schema is itself Validable
type Schema map[FieldTag]Validable

var _ Validable = (*Schema)(nil)

// A Souuup Instance
type Souuup struct {
	state  *souuupState
	schema Schema
}

func NewSouuup(schema Schema) *Souuup {
	return &Souuup{
		state:  &souuupState{NewValidationError()},
		schema: schema,
	}
}

func (u *Souuup) Validate() error {
	u.schema.Validate(u.state.errors, "")

	if u.state.errors.HasErrors() {
		return u.state.errors
	}
	return nil
}

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

func (s Schema) Errors() *ValidationError {
	errors := NewValidationError()
	s.Validate(errors, "")
	return errors
}

type Validable interface {
	Validate(*ValidationError, FieldTag)
	Errors() *ValidationError
}
