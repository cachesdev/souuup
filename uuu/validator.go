// main package of Soup
package u

// FieldTag represents the "key" of a field, and will be used to identify a field on
// an error map
type FieldTag = string

// Internal State
type souuupState struct {
	errors *ValidationError
}

// Schema of all validable fields/nested fields
// Schema is itself Validable
type Schema map[FieldTag]Validable

var _ Validable = (*Schema)(nil)

// A Souup Instance
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
	u.schema.Validate(u.state.errors)

	if u.state.errors.HasErrors() {
		return u.state.errors
	}
	return nil
}

func (d Schema) Validate(errors *ValidationError) {
	for tag, field := range d {
		errors.NestedErrors[tag] = NewValidationError()
		errors.NestedErrors[tag].Parent = errors

		field.Validate(errors.NestedErrors[tag])
	}
}

func (d Schema) Errors() *ValidationError {
	return d.Errors()
}

type Validable interface {
	Validate(*ValidationError)
	Errors() *ValidationError
}
