package u

// FieldState holds the value being validated and any validation errors.
// It is passed to validation rules to provide access to the value.
type FieldState[T any] struct {
	Value  T
	errors *ValidationError
}

// FieldDef represents a field with its value and validation rules.
// It implements the Validable interface, allowing it to be used in schemas.
type FieldDef[T any] struct {
	state FieldState[T]
	rules []Rule[T]
}

var _ Validable = (*FieldDef[any])(nil)

// Field is the main function for creating a validatable field. It takes a value to validate,
// and a set of rules that will match the type of the given value.
//
// Example:
//
//	// Validate a string with minimum and maximum length rules
//	nameField := u.Field("John Doe", u.MinS(2), u.MaxS(50))
//
//	// Validate a number with minimum value rule
//	ageField := u.Field(25, u.MinN(18))
//
//	// Validate using a custom rule
//	emailField := u.Field("user@example.com", func(fs u.FieldState[string]) error {
//		if !strings.Contains(fs.Value, "@") {
//			return fmt.Errorf("invalid email format")
//		}
//		return nil
//	})
func Field[T any](value T, rules ...Rule[T]) *FieldDef[T] {
	return &FieldDef[T]{
		state: FieldState[T]{
			Value: value,
		},
		rules: rules,
	}
}

// Validate applies all rules to the field and adds any validation errors to the provided
// ValidationError object under the specified tag.
func (f *FieldDef[T]) Validate(ve *ValidationError, tag FieldTag) {
	for _, rule := range f.rules {
		ruleErr := rule(f.state)
		if ruleErr != nil {
			ve.AddError(tag, ruleErr)
		}
	}
}

// Errors returns the validation errors associated with this field.
func (f FieldDef[T]) Errors() *ValidationError {
	return f.state.errors
}
