package u

type FieldState[T any] struct {
	value  T
	errors *ValidationError
}

type FieldDef[T any] struct {
	state FieldState[T]
	rules []Rule[T]
}

var _ Validable = (*FieldDef[any])(nil)

// Field is the main function in Validator. It takes a value to validate,
// and a set of rules that will match the type of the given value.
func Field[T any](value T, rules ...Rule[T]) *FieldDef[T] {
	return &FieldDef[T]{
		state: FieldState[T]{
			value: value,
		},
		rules: rules,
	}
}

func (f *FieldDef[T]) Validate(ve *ValidationError, tag FieldTag) {
	for _, rule := range f.rules {
		ruleErr := rule(f.state)
		if ruleErr != nil {
			ve.AddError(tag, ruleErr)
		}
	}
}

func (f FieldDef[T]) Errors() *ValidationError {
	return f.state.errors
}
