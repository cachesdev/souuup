package u_test

import (
	"errors"
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
