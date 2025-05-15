package u

type FieldState[T any] struct {
	value  T
	tag    string
	errors *ValidationError
}

type FieldDef[T any] struct {
	state FieldState[T]
	rules []Rule[T]
}

var _ Validable = (*FieldDef[any])(nil)

func (f *FieldDef[T]) Validate(errors *ValidationError) {
	for _, rule := range f.rules {
		ruleErr := rule(f.state)
		if ruleErr != nil {
			errors.AddError(f.state.tag, ruleErr)
		}
	}
}

// Returns a FieldTag
func (f FieldDef[T]) Tag() string {
	return f.state.tag
}

// Sets a FieldTag
func (f FieldDef[T]) SetTag(tag string) {
	f.state.tag = tag
}

func (f FieldDef[T]) Errors() *ValidationError {
	return f.state.errors
}

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
