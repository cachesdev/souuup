// main package of Soup
package u

// FieldTag represents the "key" of a field, and will be used to identify a field on
// an error map
type FieldTag = string

// a Souuup instance
type Souuup map[FieldTag]Validable

// a Souuup instance for boring people
type Validator = Souuup

type souuupState struct {
	errors *ValidationError
}

// Validate runs validation on all fields and returns an error if any fail
func (v Souuup) Validate() error {
	errors := NewValidationError()
	hasErrors := false

	for tag, field := range v {
		field.SetTag(tag)

		// Run validation on the field
		if err := field.Validate(); err != nil {
			hasErrors = true
			errors.AddError(tag, field.Errors())
		}
	}

	if hasErrors {
		return errors
	}
	return nil
}

// VeFromError casts an error to ValidationError if possible
func VeFromError(err error) (*ValidationError, bool) {
	if err == nil {
		return nil, false
	}

	ve, ok := err.(*ValidationError)
	return ve, ok
}

func (v Souuup) Errors() *ValidationError {
	errors := NewValidationError()
	for tag, field := range v {
		fieldErrors := field.Errors()
		if fieldErrors != nil && fieldErrors.HasErrors() {
			errors.NestedErrors[tag] = fieldErrors
		}
	}
	return errors
}

type Validable interface {
	Validate() error
	Tag() string
	SetTag(string)
	Errors() *ValidationError
}
