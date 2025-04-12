package u_test

import (
	"testing"

	u "github.com/cachesdev/souuup/uuu"
)

func TestValidationError(t *testing.T) {
	vErr := u.NewValidationError()

	vErr.Errors = u.FieldErrors{
		"field1":       []u.RuleError{"bad", "bad, twice"},
		"field2":       []u.RuleError{"2bad", "2bad, twice"},
		"nestedField1": []u.RuleError{"this is nested", "this is nested, twice"},
	}

	vErr.NestedErrors["nestedField1"] = u.NewValidationError()
	vErr.NestedErrors["nestedField1"].Errors = u.FieldErrors{
		"nestedField1Property": []u.RuleError{"invalid format", "too short"},
	}

	vErr.NestedErrors["nestedField1"].NestedErrors["deeplyNested"] = u.NewValidationError()
	vErr.NestedErrors["nestedField1"].NestedErrors["deeplyNested"].Errors = u.FieldErrors{
		"deepProperty": []u.RuleError{"cannot be null", "value out of range"},
	}

	vErr.NestedErrors["nestedField2"] = u.NewValidationError()
	vErr.NestedErrors["nestedField2"].Errors = u.FieldErrors{
		"someOtherProperty": []u.RuleError{"failed validation", "requires attention"},
	}

	if !vErr.HasErrors() {
		t.Error("Expected validation error to have errors")
	}
}

func TestValidationError2(t *testing.T) {
	vErr := u.NewValidationError()

	vErr.Errors = u.FieldErrors{
		"field1":        u.RuleErrors{"a validation error"},
		"nestedStruct1": u.RuleErrors{"a validation error"},
	}

	vErr.NestedErrors["nestedStruct1"] = u.NewValidationError()
	vErr.NestedErrors["nestedStruct1"].Errors = u.FieldErrors{
		"nestedField1": u.RuleErrors{"a validation error"},
	}

	if !vErr.HasErrors() {
		t.Error("Expected validation error to have errors")
	}
}
