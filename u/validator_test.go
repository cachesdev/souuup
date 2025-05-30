package u_test

import (
	"strings"
	"testing"

	"github.com/cachesdev/souuup/u"
)

func TestSchema_Validate(t *testing.T) {
	t.Run("empty schema produces no errors", func(t *testing.T) {
		// Arrange
		schema := u.Schema{}
		ve := u.NewValidationError()

		// Act
		schema.Validate(ve, "")

		// Assert
		if ve.HasErrors() {
			t.Error("expected no validation errors for empty schema")
		}
	})

	t.Run("schema with valid fields produces no errors", func(t *testing.T) {
		// Arrange
		schema := u.Schema{
			"field1": &mockValidable{hasErrors: false},
			"field2": &mockValidable{hasErrors: false},
		}
		ve := u.NewValidationError()

		// Act
		schema.Validate(ve, "")

		// Assert
		if ve.HasErrors() {
			t.Error("expected no validation errors for schema with valid fields")
		}
	})

	t.Run("schema with invalid fields produces errors", func(t *testing.T) {
		// Arrange
		schema := u.Schema{
			"field1": &mockValidable{hasErrors: true},
			"field2": &mockValidable{hasErrors: false},
		}
		ve := u.NewValidationError()

		// Act
		schema.Validate(ve, "")

		// Assert
		if !ve.HasErrors() {
			t.Error("expected validation errors for schema with invalid fields")
		}
	})

	t.Run("validates all fields in schema", func(t *testing.T) {
		// Arrange
		field1 := &mockValidable{hasErrors: false}
		field2 := &mockValidable{hasErrors: false}
		field3 := &mockValidable{hasErrors: false}

		schema := u.Schema{
			"field1": field1,
			"field2": field2,
			"field3": field3,
		}
		ve := u.NewValidationError()

		// Act
		schema.Validate(ve, "")

		// Assert
		assertFieldValidated(t, field1, "field1")
		assertFieldValidated(t, field2, "field2")
		assertFieldValidated(t, field3, "field3")
	})

	t.Run("validates nested schemas correctly", func(t *testing.T) {
		// Arrange
		nestedField := &mockValidable{hasErrors: true}

		schema := u.Schema{
			"root": u.Schema{
				"nestedField": nestedField,
			},
		}
		ve := u.NewValidationError()

		// Act
		schema.Validate(ve, "")

		// Assert
		if !ve.HasErrors() {
			t.Error("expected validation errors from nested schema")
		}

		if _, exists := ve.NestedErrors["root"]; !exists {
			t.Error("expected nested error for 'root' field")
		}

		if !nestedField.validateCalled {
			t.Error("expected Validate() to be called on nested field")
		}
	})

	t.Run("deeply nested schemas are validated correctly", func(t *testing.T) {
		// Arrange
		deepField := &mockValidable{hasErrors: true}

		rootSchema := u.Schema{
			"root": u.Schema{
				"middle": u.Schema{
					"deepField": deepField,
				},
			},
		}
		ve := u.NewValidationError()

		// Act
		rootSchema.Validate(ve, "")

		// Assert
		if !ve.HasErrors() {
			t.Error("expected validation errors from deeply nested schema")
		}

		if !deepField.validateCalled {
			t.Error("expected Validate() to be called on deeply nested field")
		}

		if errors, exists := ve.NestedErrors["root"]; !exists {
			t.Error("expected nested validation error for 'root' field")
		} else {
			if !errors.HasErrors() {
				t.Error("expected nested errors for 'root' field")
			}
		}

		rootErrors := ve.NestedErrors["root"]
		if errors, exists := rootErrors.NestedErrors["middle"]; !exists {
			t.Error("expected nested validation error for 'middle' field")
		} else {
			if !errors.HasErrors() {
				t.Error("expected nested errors for 'middle' field")
			}
		}
	})

	t.Run("deeply nested schemas are validated correctly when valid", func(t *testing.T) {
		// Arrange
		deepField := &mockValidable{hasErrors: false}

		rootSchema := u.Schema{
			"root": u.Schema{
				"middle": u.Schema{
					"deepField": deepField,
				},
			},
		}
		ve := u.NewValidationError()

		// Act
		rootSchema.Validate(ve, "")

		// Assert
		if ve.HasErrors() {
			t.Error("expected no validation errors from deeply nested schema")
		}

		if !deepField.validateCalled {
			t.Error("expected Validate() to be called on deeply nested field")
		}

		if errors, exists := ve.NestedErrors["root"]; exists {
			if errors.HasErrors() {
				t.Error("expected no nested error for 'root' field")
			}
		}

		rootErrors := ve.NestedErrors["root"]
		if errors, exists := rootErrors.NestedErrors["middle"]; exists {
			if errors.HasErrors() {
				t.Error("expected no nested error for 'middle' field")
			}
		}
	})

	t.Run("mixed direct and nested fields validate correctly", func(t *testing.T) {
		// Arrange
		directField := &mockValidable{hasErrors: true}
		nestedField := &mockValidable{hasErrors: true}

		schema := u.Schema{
			"direct": directField,
			"nested": u.Schema{
				"nestedField": nestedField,
			},
		}
		ve := u.NewValidationError()

		// Act
		schema.Validate(ve, "")

		// Assert
		if !ve.HasErrors() {
			t.Error("expected validation errors from mixed schema")
		}

		if !directField.validateCalled {
			t.Error("expected Validate() to be called on direct field")
		}

		if !nestedField.validateCalled {
			t.Error("expected Validate() to be called on nested field")
		}

		if _, exists := ve.NestedErrors["nested"]; !exists {
			t.Error("expected nested error for 'nested' field")
		}
	})

	t.Run("complex nested structure is validated correctly", func(t *testing.T) {
		// Arrange
		tc := newComplexSchemaTestCase()
		ve := u.NewValidationError()

		// Act
		tc.schema.Validate(ve, "")

		// Assert
		if !ve.HasErrors() {
			t.Error("expected validation errors from complex schema")
		}

		// Helper to check if a nested path has errors
		checkNestedPath := func(path string) bool {
			parts := strings.Split(path, ".")
			currentVE := ve

			for i, part := range parts {
				if i == len(parts)-1 {
					// Final part, check in Errors
					if _, exists := currentVE.Errors[part]; exists {
						return true
					}
					return false
				} else {
					// Intermediate part, check in NestedErrors
					if nestedVE, exists := currentVE.NestedErrors[part]; exists {
						currentVE = nestedVE
					} else {
						return false
					}
				}
			}
			return false
		}

		for _, path := range tc.expectedPaths {
			if !checkNestedPath(path) {
				t.Errorf("expected errors for path %q, but none found", path)
			}
		}

		for _, path := range tc.unexpectedPaths {
			if checkNestedPath(path) {
				t.Errorf("unexpected error for path %q", path)
			}
		}
	})
}

func TestSchema_Errors(t *testing.T) {
	t.Run("returns ValidationError with no errors for valid schema", func(t *testing.T) {
		// Arrange
		schema := u.Schema{
			"field1": &mockValidable{hasErrors: false},
			"field2": &mockValidable{hasErrors: false},
		}

		// Act
		errors := schema.Errors()

		// Assert
		if errors == nil {
			t.Error("expected non-nil ValidationError")
		}

		if errors.HasErrors() {
			t.Error("expected no validation errors for valid schema")
		}
	})

	t.Run("returns ValidationError with errors for invalid schema", func(t *testing.T) {
		// Arrange
		schema := u.Schema{
			"field1": &mockValidable{hasErrors: true},
			"field2": &mockValidable{hasErrors: false},
		}

		// Act
		errors := schema.Errors()

		// Assert
		if !errors.HasErrors() {
			t.Error("expected validation errors for invalid schema")
		}
	})
}

func TestSouuup_Validate(t *testing.T) {
	t.Run("returns nil for valid empty schema", func(t *testing.T) {
		// Arrange
		schema := u.Schema{}
		souuup := u.NewSouuup(schema)

		// Act
		err := souuup.Validate()

		// Assert
		if err != nil {
			t.Errorf("expected nil error for valid empty schema, got %v", err)
		}
	})

	t.Run("returns nil for valid schema with fields", func(t *testing.T) {
		// Arrange
		schema := u.Schema{
			"field1": &mockValidable{hasErrors: false},
			"field2": &mockValidable{hasErrors: false},
		}
		souuup := u.NewSouuup(schema)

		// Act
		err := souuup.Validate()

		// Assert
		if err != nil {
			t.Errorf("expected nil error for valid schema with fields, got %v", err)
		}
	})

	t.Run("returns validation error for invalid schema", func(t *testing.T) {
		// Arrange
		schema := u.Schema{
			"field1": &mockValidable{hasErrors: true, errorMessage: "field1 error"},
		}
		souuup := u.NewSouuup(schema)

		// Act
		err := souuup.Validate()

		// Assert
		if err == nil {
			t.Error("expected validation error for invalid schema, got nil")
		}

		if !strings.Contains(err.Error(), "field1 error") {
			t.Errorf("error message %q does not contain expected error message", err.Error())
		}
	})

	t.Run("validates nested fields correctly", func(t *testing.T) {
		// Arrange
		nestedField := &mockValidable{hasErrors: true, errorMessage: "nested error"}
		schema := u.Schema{
			"parent": u.Schema{
				"nested": nestedField,
			},
		}
		souuup := u.NewSouuup(schema)

		// Act
		err := souuup.Validate()

		// Assert
		if err == nil {
			t.Error("expected validation error for schema with invalid nested field")
		}

		if !nestedField.validateCalled {
			t.Error("expected Validate() to be called on nested field")
		}

		if !strings.Contains(err.Error(), "nested error") {
			t.Errorf("error message %q does not contain expected nested error message", err.Error())
		}
	})

	t.Run("handles complex mixed valid/invalid schema", func(t *testing.T) {
		// Arrange
		validField := &mockValidable{hasErrors: false}
		invalidField1 := &mockValidable{hasErrors: true, errorMessage: "error1"}
		invalidField2 := &mockValidable{hasErrors: true, errorMessage: "error2"}

		schema := u.Schema{
			"valid":    validField,
			"invalid1": invalidField1,
			"parent": u.Schema{
				"invalid2": invalidField2,
				"valid2":   &mockValidable{hasErrors: false},
			},
		}
		souuup := u.NewSouuup(schema)

		// Act
		err := souuup.Validate()

		// Assert
		if err == nil {
			t.Error("expected validation error for mixed valid/invalid schema")
		}

		if !validField.validateCalled {
			t.Error("expected Validate() to be called on valid field")
		}
		if !invalidField1.validateCalled {
			t.Error("expected Validate() to be called on invalid field 1")
		}
		if !invalidField2.validateCalled {
			t.Error("expected Validate() to be called on invalid field 2")
		}

		errStr := err.Error()
		if !strings.Contains(errStr, "error1") {
			t.Errorf("error message %q does not contain expected error message 'error1'", errStr)
		}
		if !strings.Contains(errStr, "error2") {
			t.Errorf("error message %q does not contain expected error message 'error2'", errStr)
		}
	})

	t.Run("complex schema error paths are correctly formed", func(t *testing.T) {
		// Arrange
		tc := newComplexTestCase()
		souuup := u.NewSouuup(tc.schema)

		// Act
		err := souuup.Validate()

		// Assert
		if err == nil {
			t.Error("expected validation error for complex schema")
		}

		errStr := err.Error()
		for _, errMsg := range tc.expectedErrorMessages {
			if !strings.Contains(errStr, errMsg) {
				t.Errorf("error message %q does not contain expected message %q", errStr, errMsg)
			}
		}
	})
}

func TestValidator_NewSouuup(t *testing.T) {
	t.Run("returns a pointer to Souuup", func(t *testing.T) {
		// Act
		schema := u.Schema{}
		s := u.NewSouuup(schema)

		// Assert
		if s == nil {
			t.Error("expected NewSouuup() to return a non-nil pointer")
		}
	})

	t.Run("initialises with empty schema", func(t *testing.T) {
		// Arrange
		schema := u.Schema{}

		// Act
		s := u.NewSouuup(schema)

		// Assert
		err := s.Validate()
		if err != nil {
			t.Errorf("expected no validation errors for empty schema, got %v", err)
		}
	})

	t.Run("initialises with provided schema", func(t *testing.T) {
		// Arrange
		schema := u.Schema{
			"test": &mockValidable{
				hasErrors: false,
			},
		}

		// Act
		s := u.NewSouuup(schema)

		// Assert
		err := s.Validate()
		if err != nil {
			t.Errorf("expected no validation errors for valid schema, got %v", err)
		}
	})

	t.Run("validates schema on validate call", func(t *testing.T) {
		// Arrange
		mockField := &mockValidable{
			hasErrors: true,
		}
		schema := u.Schema{
			"test": mockField,
		}

		// Act
		s := u.NewSouuup(schema)
		err := s.Validate()

		// Assert
		if err == nil {
			t.Error("expected validation errors but got nil")
		}
		if !mockField.validateCalled {
			t.Error("expected Validate() to be called on schema field")
		}
	})
}

// Mock implementation of Validable interface for testing
type mockValidable struct {
	validateCalled bool
	hasErrors      bool
	errorMessage   string
}

func (m *mockValidable) Validate(ve *u.ValidationError, tag u.FieldTag) {
	m.validateCalled = true
	if m.hasErrors {
		errorMsg := "mock validation error"
		if m.errorMessage != "" {
			errorMsg = m.errorMessage
		}
		ve.AddError(tag, u.RuleError(errorMsg))
	}
}

func (m *mockValidable) Errors() *u.ValidationError {
	ve := u.NewValidationError()
	if m.hasErrors {
		errorMsg := "mock validation error"
		if m.errorMessage != "" {
			errorMsg = m.errorMessage
		}
		ve.AddError("mock", u.RuleError(errorMsg))
	}
	return ve
}

// Helper assertion functions
func assertFieldValidated(t *testing.T, field *mockValidable, fieldName string) {
	t.Helper()
	if !field.validateCalled {
		t.Errorf("expected Validate() to be called on field %q", fieldName)
	}
}

type complexSchemaTestCase struct {
	schema          u.Schema
	expectedPaths   []string // Paths to fields that should have errors
	unexpectedPaths []string // Paths to fields that should not have errors
}

func newComplexSchemaTestCase() complexSchemaTestCase {
	userSchema := u.Schema{
		"name":  &mockValidable{hasErrors: true, errorMessage: "invalid name"},
		"email": &mockValidable{hasErrors: false},
		"address": u.Schema{
			"street": &mockValidable{hasErrors: true, errorMessage: "invalid street"},
			"city":   &mockValidable{hasErrors: false},
		},
	}

	billingSchema := u.Schema{
		"name": &mockValidable{hasErrors: false},
		"card": &mockValidable{hasErrors: true, errorMessage: "invalid card"},
	}

	rootSchema := u.Schema{
		"user":    userSchema,
		"billing": billingSchema,
		"direct":  &mockValidable{hasErrors: true, errorMessage: "direct error"},
	}

	return complexSchemaTestCase{
		schema: rootSchema,
		expectedPaths: []string{
			"user.name",
			"user.address.street",
			"billing.card",
			"direct",
		},
		unexpectedPaths: []string{
			"user.email",
			"user.address.city",
			"billing.name",
		},
	}
}

type complexTestCase struct {
	schema                u.Schema
	expectedErrorMessages []string
}

func newComplexTestCase() complexTestCase {
	return complexTestCase{
		schema: u.Schema{
			"username": &mockValidable{hasErrors: true, errorMessage: "username too short"},
			"email":    &mockValidable{hasErrors: true, errorMessage: "invalid email format"},
			"profile": u.Schema{
				"bio": &mockValidable{hasErrors: true, errorMessage: "bio too long"},
				"age": &mockValidable{hasErrors: false},
			},
			"settings": u.Schema{
				"notifications": u.Schema{
					"email": &mockValidable{hasErrors: true, errorMessage: "invalid notification setting"},
				},
			},
		},
		expectedErrorMessages: []string{
			"username too short",
			"invalid email format",
			"bio too long",
			"invalid notification setting",
		},
	}
}
