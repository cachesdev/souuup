package u_test

import (
	"errors"
	"testing"

	"github.com/cachesdev/souuup/u"
)

func TestField_Field(t *testing.T) {
	t.Run("returns a pointer to FieldDef", func(t *testing.T) {
		// Act
		field := u.Field("test")

		// Assert
		if field == nil {
			t.Error("expected Field() to return a non-nil pointer")
		}
	})

	t.Run("initialises with empty errors", func(t *testing.T) {
		// Act
		field := u.Field("test")

		// Assert
		if field.Errors() != nil {
			t.Error("expected field.Errors() to be nil for newly created Field")
		}
	})
}

func TestField_Validate(t *testing.T) {
	// Helper function to create a rule that always returns an error
	errorRule := func(msg string) u.Rule[int] {
		return func(state u.FieldState[int]) error {
			return errors.New(msg)
		}
	}

	// Helper function to create a rule that always passes validation
	passRule := func() u.Rule[int] {
		return func(state u.FieldState[int]) error {
			return nil
		}
	}

	tests := []struct {
		name          string
		value         int
		rules         []u.Rule[int]
		tag           u.FieldTag
		expectedError bool
		errorField    string
		errorMessages u.RuleErrors
	}{
		{
			name:          "passes validation with no rules",
			value:         123,
			rules:         []u.Rule[int]{},
			tag:           "field",
			expectedError: false,
		},
		{
			name:          "passes validation with passing rule",
			value:         123,
			rules:         []u.Rule[int]{passRule()},
			tag:           "field",
			expectedError: false,
		},
		{
			name:          "fails validation with one failing rule",
			value:         123,
			rules:         []u.Rule[int]{errorRule("validation failed")},
			tag:           "username",
			expectedError: true,
			errorField:    "username",
			errorMessages: u.RuleErrors{"validation failed"},
		},
		{
			name:          "reports multiple errors from multiple rules",
			value:         123,
			rules:         []u.Rule[int]{errorRule("error 1"), errorRule("error 2")},
			tag:           "password",
			expectedError: true,
			errorField:    "password",
			errorMessages: u.RuleErrors{"error 1", "error 2"},
		},
		{
			name:          "reports only errors from failing rules",
			value:         123,
			rules:         []u.Rule[int]{passRule(), errorRule("only error"), passRule()},
			tag:           "email",
			expectedError: true,
			errorField:    "email",
			errorMessages: u.RuleErrors{"only error"},
		},
		{
			name:          "uses provided field tag for error",
			value:         123,
			rules:         []u.Rule[int]{errorRule("validation error")},
			tag:           "custom_field",
			expectedError: true,
			errorField:    "custom_field",
			errorMessages: u.RuleErrors{"validation error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			field := u.Field(tc.value, tc.rules...)
			ve := u.NewValidationError()

			// Act
			field.Validate(ve, tc.tag)

			// Assert
			hasErrors := ve.HasErrors()
			if tc.expectedError != hasErrors {
				t.Errorf("expected HasErrors() to be %v, got %v", tc.expectedError, hasErrors)
			}

			if !tc.expectedError {
				// No need to check error details if no error was expected
				return
			}

			// Check that the specific field has errors
			fieldErrors, found := ve.Errors[tc.errorField]
			if !found {
				t.Errorf("expected errors for field %q, but found none", tc.errorField)
			}

			// Check that no other fields have errors
			for field := range ve.Errors {
				if field != tc.errorField {
					t.Errorf("expected no errors for field %q, but found %d", field, len(ve.Errors[field]))
				}
			}

			// Check error count
			if len(fieldErrors) != len(tc.errorMessages) {
				t.Errorf("expected %d error(s) for field %q, got %d",
					len(tc.errorMessages), tc.errorField, len(fieldErrors))
			}

			// Check error messages
			for i, expected := range tc.errorMessages {
				if i >= len(fieldErrors) {
					t.Errorf("missing expected error: %q", expected)
					continue
				}
				if fieldErrors[i] != expected {
					t.Errorf("expected error message %q, got %q", expected, fieldErrors[i].Error())
				}
			}
		})
	}
}

func TestField_Errors(t *testing.T) {
	t.Run("returns nil for a newly created field", func(t *testing.T) {
		// Act
		field := u.Field(123)

		// Assert
		if field.Errors() != nil {
			t.Error("expected Errors() to return nil for a newly created field")
		}
	})

	t.Run("returns error object after validation with errors", func(t *testing.T) {
		// Arrange
		errorRule := func(state u.FieldState[int]) error {
			return errors.New("validation error")
		}
		field := u.Field(123, errorRule)
		ve := u.NewValidationError()

		// Act
		field.Validate(ve, "test_field")

		// Assert
		if !ve.HasErrors() {
			t.Error("expected validation to fail")
		}

		fieldErrors := ve.Errors["test_field"]
		if len(fieldErrors) != 1 {
			t.Errorf("expected 1 error, got %d", len(fieldErrors))
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		// Act
		field1 := u.Field("string value")
		field2 := u.Field(42)
		field3 := u.Field(true)
		field4 := u.Field([]string{"slice", "of", "strings"})

		// Assert
		if field1.Errors() != nil || field2.Errors() != nil ||
			field3.Errors() != nil || field4.Errors() != nil {
			t.Error("expected Errors() to return nil for all newly created fields regardless of type")
		}
	})
}
