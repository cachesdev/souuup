package u

import (
	"fmt"
	"maps"
)

type FieldState[T any] struct {
	value  T
	tag    string
	errors *ValidationError
}

type FieldDef[T any] struct {
	state FieldState[T]
	rules []Rule[T]
}

func (f *FieldDef[T]) Validate() error {
	errors := NewValidationError()
	hasErrors := false

	for _, rule := range f.rules {
		ruleErr := rule(f.state)
		if ruleErr != nil {
			hasErrors = true

			if validErr, ok := VeFromError(ruleErr); ok {
				maps.Copy(errors.NestedErrors, validErr.NestedErrors)

				for fieldTag, fieldErrors := range validErr.Errors {
					for _, errMsg := range fieldErrors {
						errors.AddError(fieldTag, errMsg)
					}
				}
			} else {
				fmt.Println("erm")
				errors.AddError(f.state.tag, ruleErr)
			}
		}
	}

	if hasErrors {
		return errors
	}
	return nil
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

func NestedFn[T any](value T, nestedFn func(T) Souuup) *FieldDef[T] {
	return &FieldDef[T]{
		state: FieldState[T]{
			value: value,
		},
		rules: []Rule[T]{
			func(state FieldState[T]) error {
				nestedValidator := nestedFn(state.value)

				if err := nestedValidator.ValidateSouuup(); err != nil {
					return err
				}

				return nil
			},
		},
	}
}

func Nested(uuu Souuup) *FieldDef[struct{}] {
	return &FieldDef[struct{}]{
		state: FieldState[struct{}]{
			value: struct{}{},
		},
		rules: []Rule[struct{}]{
			func(state FieldState[struct{}]) error {
				if err := uuu.ValidateSouuup(); err != nil {
					return err
				}
				return nil
			},
		},
	}
}
