package u

type FieldState[T any] struct {
	value  T
	tag    string
	errors []error
}

type FieldDef[T any] struct {
	state FieldState[T]
	rules []Rule[T]
}

// A single implementation for all types
func (f FieldDef[T]) Validate() bool {
	for _, rule := range f.rules {
		if !rule(f.state) {
			return false
		}
	}
	return true
}

// Returns a FieldTag
func (f FieldDef[T]) Tag() string {
	return f.state.tag
}

// Sets a FieldTag
func (f FieldDef[T]) SetTag(tag string) {
	f.state.tag = tag
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
