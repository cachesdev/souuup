package u_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"

	u "github.com/cachesdev/souuup/uuu"
)

func TestValidationError_NewValidationError(t *testing.T) {
	t.Run("returns a pointer to ValidationError", func(t *testing.T) {
		// Act
		ve := u.NewValidationError()

		// Assert
		if ve == nil {
			t.Error("expected NewValidationError() to return a non-nil pointer")
		}
	})

	t.Run("initialises with empty fields", func(t *testing.T) {
		// Act
		ve := u.NewValidationError()

		// Assert
		if ve.Errors == nil {
			t.Error("expected Errors map to be initialised, but it was nil")
		}
		assertErroredFieldsLen(t, ve, 0)
		if ve.NestedErrors == nil {
			t.Error("expected NestedErrors map to be initialised, but it was nil")
		}
		if len(ve.NestedErrors) != 0 {
			t.Errorf("expected NestedErrors map to be empty, but it had %d entries", len(ve.NestedErrors))
		}
		if ve.Parent != nil {
			t.Error("expected Parent to be nil for newly created ValidationError")
		}
	})
}

func TestValidationError_AddError(t *testing.T) {
	t.Run("adds error to empty ValidationError", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()
		testField := "username"
		testErrMsg := "must be at least 3 characters"
		testErr := errors.New(testErrMsg)

		// Act
		ve.AddError(testField, testErr)

		// Assert
		assertErroredFieldsLen(t, ve, 1)
		assertFieldHasErrors(t, testField, ve.Errors[testField])
		assertErrorsLen(t, testField, ve.Errors[testField], 1)
		assertErrorMessage(t, ve.Errors[testField][0].Error(), testErrMsg)
	})

	t.Run("appends error to existing field errors", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()
		testField := "email"
		firstErrMsg := "cannot be empty"
		firstErr := errors.New(firstErrMsg)
		secondErrMsg := "must be a valid email address"
		secondErr := errors.New(secondErrMsg)

		// Act
		ve.AddError(testField, firstErr)
		ve.AddError(testField, secondErr)

		// Assert
		assertErroredFieldsLen(t, ve, 1)
		assertFieldHasErrors(t, testField, ve.Errors[testField])
		assertErrorsLen(t, testField, ve.Errors[testField], 2)
		assertErrorMessage(t, ve.Errors[testField][0].Error(), firstErrMsg)
		assertErrorMessage(t, ve.Errors[testField][1].Error(), secondErrMsg)
	})

	t.Run("adds errors to multiple fields", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()
		usernameField := "username"
		usernameErrMsg := "must be at least 3 characters"
		usernameErr := errors.New(usernameErrMsg)
		emailField := "email"
		emailErrMsg := "must be a valid email address"
		emailErr := errors.New(emailErrMsg)

		// Act
		ve.AddError(usernameField, usernameErr)
		ve.AddError(emailField, emailErr)

		// Assert
		assertErroredFieldsLen(t, ve, 2)
		assertFieldHasErrors(t, usernameField, ve.Errors[usernameField])
		assertFieldHasErrors(t, emailField, ve.Errors[emailField])
		assertErrorMessage(t, ve.Errors[usernameField][0].Error(), usernameErrMsg)
		assertErrorMessage(t, ve.Errors[emailField][0].Error(), emailErrMsg)
	})
}

func TestValidationError_HasErrors(t *testing.T) {
	t.Run("returns false for empty ValidationError", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()

		// Act
		hasErrors := ve.HasErrors()

		// Assert
		assertHasErrorsFalse(t, hasErrors)
	})

	t.Run("returns true when direct errors exist", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()
		testField := "username"
		testErrMsg := "must be at least 3 characters"
		testErr := errors.New(testErrMsg)
		ve.AddError(testField, testErr)

		// Act
		hasErrors := ve.HasErrors()

		// Assert
		assertHasErrorsTrue(t, hasErrors)
	})

	t.Run("returns true when nested errors exist", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()
		nestedField := "address"
		nested := ve.GetOrCreateNested(nestedField)
		testField := "street"
		testErrMsg := "cannot be empty"
		testErr := errors.New(testErrMsg)
		nested.AddError(testField, testErr)

		// Act
		hasErrors := ve.HasErrors()

		// Assert
		assertHasErrorsTrue(t, hasErrors)
	})

	t.Run("returns true when deeply nested errors exist", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()
		level1Field := "address"
		level1 := ve.GetOrCreateNested(level1Field)
		level2Field := "billing"
		level2 := level1.GetOrCreateNested(level2Field)
		testField := "postcode"
		testErrMsg := "invalid format"
		testErr := errors.New(testErrMsg)
		level2.AddError(testField, testErr)

		// Act
		hasErrors := ve.HasErrors()

		// Assert
		assertHasErrorsTrue(t, hasErrors)
	})

	t.Run("returns false when nested ValidationErrors have no errors", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()
		// Create nested ValidationError but don't add any errors to it
		ve.GetOrCreateNested("address")

		// Act
		hasErrors := ve.HasErrors()

		// Assert
		assertHasErrorsFalse(t, hasErrors)
	})
}

func TestValidationError_ToMap(t *testing.T) {
	t.Run("returns nil for ValidationError with no errors", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()

		// Act
		result := ve.ToMap()

		// Assert
		if result != nil {
			t.Errorf("expected nil for ValidationError with no errors, got %v", result)
		}
	})

	verifyMap := func(t *testing.T, result, expected u.ToMapResult) {
		t.Helper()

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("map mismatch:expected %v, got %v", expected, result)
		}
	}

	tests := []struct {
		name     string
		setup    func() *u.ValidationError
		expected u.ToMapResult
	}{
		{
			name: "direct field errors only",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("username", errors.New("must be at least 3 characters"))
				ve.AddError("email", errors.New("cannot be empty"))
				return ve
			},
			expected: u.ToMapResult{
				"username": {"errors": u.RuleErrors{"must be at least 3 characters"}},
				"email":    {"errors": u.RuleErrors{"cannot be empty"}},
			},
		},
		{
			name: "nested field errors only (one level)",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				nested := ve.GetOrCreateNested("address")
				nested.AddError("street", errors.New("cannot be empty"))
				nested.AddError("city", errors.New("cannot be empty"))
				return ve
			},
			expected: u.ToMapResult{
				"address": {
					"street": map[string]any{"errors": u.RuleErrors{"cannot be empty"}},
					"city":   map[string]any{"errors": u.RuleErrors{"cannot be empty"}},
				},
			},
		},
		{
			name: "combination of direct and nested errors",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("username", errors.New("must be at least 3 characters"))
				nested := ve.GetOrCreateNested("address")
				nested.AddError("street", errors.New("cannot be empty"))
				return ve
			},
			expected: u.ToMapResult{
				"username": {
					"errors": u.RuleErrors{"must be at least 3 characters"},
				},
				"address": {
					"street": map[string]any{"errors": u.RuleErrors{"cannot be empty"}},
				},
			},
		},
		{
			name: "deeply nested errors (multiple levels)",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				level1 := ve.GetOrCreateNested("user")
				level2 := level1.GetOrCreateNested("address")
				level2.AddError("postcode", errors.New("invalid format"))
				return ve
			},
			expected: u.ToMapResult{
				"user": {
					"address": map[string]any{
						"postcode": map[string]any{"errors": u.RuleErrors{"invalid format"}},
					},
				},
			},
		},
		{
			name: "field with both direct and nested errors",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("address", errors.New("invalid address"))
				nested := ve.GetOrCreateNested("address")
				nested.AddError("street", errors.New("cannot be empty"))
				return ve
			},
			expected: u.ToMapResult{
				"address": {
					"errors": u.RuleErrors{"invalid address"},
					"street": map[string]any{"errors": u.RuleErrors{"cannot be empty"}},
				},
			},
		},
		{
			name: "multiple errors for the same field",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("password", errors.New("too short"))
				ve.AddError("password", errors.New("needs special characters"))
				ve.AddError("password", errors.New("needs numbers"))
				return ve
			},
			expected: u.ToMapResult{
				"password": {"errors": u.RuleErrors{
					"too short",
					"needs special characters",
					"needs numbers",
				}},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ve := tc.setup()

			// Act
			result := ve.ToMap()

			// Assert
			verifyMap(t, result, tc.expected)
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	t.Run("returns empty string for ValidationError with no errors", func(t *testing.T) {
		// Arrange
		ve := u.NewValidationError()

		// Act
		result := ve.Error()

		// Assert
		if result != "" {
			t.Errorf("expected empty string for ValidationError with no errors, got %q", result)
		}
	})

	verifyErrorString := func(t *testing.T, result string, expected string, contains []string) {
		t.Helper()

		var jsonObj map[string]any
		if err := json.Unmarshal([]byte(result), &jsonObj); err != nil {
			t.Errorf("error string is not valid JSON: %v", err)
		}

		// If expected is provided, check for an exact match
		if expected != "" && result != expected {
			t.Errorf("expected exact error %q, got %q", expected, result)
			return
		}

		// Else, check for substrings being contained in the result
		for _, substr := range contains {
			if !strings.Contains(result, substr) {
				t.Errorf("expected error string to contain %q, but got %q", substr, result)
			}
		}
	}

	tests := []struct {
		name     string
		setup    func() *u.ValidationError
		expected string
		contains []string // Substrings that should be in the error message
	}{
		{
			name: "single direct error",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("username", errors.New("must be at least 3 characters"))
				return ve
			},
			expected: `{"username":{"errors":["must be at least 3 characters"]}}`,
		},
		{
			name: "multiple direct errors on different fields",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("username", errors.New("must be at least 3 characters"))
				ve.AddError("email", errors.New("must be a valid email"))
				return ve
			},
			contains: []string{
				`"username":{"errors":["must be at least 3 characters"]}`,
				`"email":{"errors":["must be a valid email"]}`,
			},
		},
		{
			name: "multiple errors on same field",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("password", errors.New("too short"))
				ve.AddError("password", errors.New("needs special characters"))
				return ve
			},
			expected: `{"password":{"errors":["too short","needs special characters"]}}`,
		},
		{
			name: "nested errors",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				nested := ve.GetOrCreateNested("address")
				nested.AddError("street", errors.New("cannot be empty"))
				return ve
			},
			expected: `{"address":{"street":{"errors":["cannot be empty"]}}}`,
		},
		{
			name: "deeply nested errors",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				level1 := ve.GetOrCreateNested("user")
				level2 := level1.GetOrCreateNested("address")
				level2.AddError("postcode", errors.New("invalid format"))
				return ve
			},
			expected: `{"user":{"address":{"postcode":{"errors":["invalid format"]}}}}`,
		},
		{
			name: "direct and nested errors combined",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("username", errors.New("invalid username"))
				nested := ve.GetOrCreateNested("profile")
				nested.AddError("bio", errors.New("too long"))
				return ve
			},
			contains: []string{
				`"username":{"errors":["invalid username"]}`,
				`"profile":{"bio":{"errors":["too long"]}}`,
			},
		},
		{
			name: "field with both direct and nested errors",
			setup: func() *u.ValidationError {
				ve := u.NewValidationError()
				ve.AddError("address", errors.New("invalid address"))
				nested := ve.GetOrCreateNested("address")
				nested.AddError("street", errors.New("cannot be empty"))
				return ve
			},
			contains: []string{
				`"address"`,
				`"errors":["invalid address"]`,
				`"street":{"errors":["cannot be empty"]}`,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ve := tc.setup()

			// Act
			result := ve.Error()

			// Assert
			verifyErrorString(t, result, tc.expected, tc.contains)
		})
	}
}

// Helpers

func assertErroredFieldsLen(t *testing.T, ve *u.ValidationError, want int) {
	t.Helper()
	if len(ve.Errors) != want {
		t.Errorf("expected %d fields with errors, got %d", want, len(ve.Errors))
	}
}

func assertFieldHasErrors(t *testing.T, fieldName string, got u.RuleErrors) {
	t.Helper()
	if len(got) == 0 {
		t.Errorf("expected errors for field %q, but found none", fieldName)
	}
}

func assertErrorsLen(t *testing.T, fieldName string, got u.RuleErrors, want int) {
	t.Helper()
	if len(got) != want {
		t.Errorf("expected %d error(s) for field %q, got %d", want, fieldName, len(got))
	}
}

func assertErrorMessage(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("expected error message %q, got %q", want, got)
	}
}

func assertHasErrorsTrue(t *testing.T, hasErrors bool) {
	t.Helper()
	if !hasErrors {
		t.Error("expected HasErrors() to return true")
	}
}

func assertHasErrorsFalse(t *testing.T, hasErrors bool) {
	t.Helper()
	if hasErrors {
		t.Error("expected HasErrors() to return false")
	}
}
