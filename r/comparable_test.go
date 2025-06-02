package r_test

import (
	"testing"

	"github.com/cachesdev/souuup/internal/testutil"
	"github.com/cachesdev/souuup/r"
	"github.com/cachesdev/souuup/u"
)

func TestNotZero(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "string: non-zero value",
			value:   "test",
			wantErr: false,
		},
		{
			name:     "string: zero value",
			value:    "",
			wantErr:  true,
			errorMsg: "value is required but has zero value",
		},
		{
			name:    "int: non-zero value",
			value:   42,
			wantErr: false,
		},
		{
			name:     "int: zero value",
			value:    0,
			wantErr:  true,
			errorMsg: "value is required but has zero value",
		},
		{
			name:    "bool: non-zero value",
			value:   true,
			wantErr: false,
		},
		{
			name:     "bool: zero value",
			value:    false,
			wantErr:  true,
			errorMsg: "value is required but has zero value",
		},
		{
			name:    "float: non-zero value",
			value:   1.23,
			wantErr: false,
		},
		{
			name:     "float: zero value",
			value:    0.0,
			wantErr:  true,
			errorMsg: "value is required but has zero value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case string:
				fs := u.FieldState[string]{Value: v}
				err := r.NotZero(fs)
				testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
			case int:
				fs := u.FieldState[int]{Value: v}
				err := r.NotZero(fs)
				testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
			case bool:
				fs := u.FieldState[bool]{Value: v}
				err := r.NotZero(fs)
				testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
			case float64:
				fs := u.FieldState[float64]{Value: v}
				err := r.NotZero(fs)
				testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
			}
		})
	}
}

func TestSameAs(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		other    any
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "string: same values",
			value:   "password",
			other:   "password",
			wantErr: false,
		},
		{
			name:     "string: different values",
			value:    "password",
			other:    "confirmPassword",
			wantErr:  true,
			errorMsg: "password does not match confirmPassword",
		},
		{
			name:    "int: same values",
			value:   42,
			other:   42,
			wantErr: false,
		},
		{
			name:     "int: different values",
			value:    42,
			other:    24,
			wantErr:  true,
			errorMsg: "42 does not match 24",
		},
		{
			name:    "bool: same values",
			value:   true,
			other:   true,
			wantErr: false,
		},
		{
			name:     "bool: different values",
			value:    true,
			other:    false,
			wantErr:  true,
			errorMsg: "true does not match false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case string:
				other := tt.other.(string)
				fs := u.FieldState[string]{Value: v}
				err := r.SameAs(other)(fs)
				testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
			case int:
				other := tt.other.(int)
				fs := u.FieldState[int]{Value: v}
				err := r.SameAs(other)(fs)
				testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
			case bool:
				other := tt.other.(bool)
				fs := u.FieldState[bool]{Value: v}
				err := r.SameAs(other)(fs)
				testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
			}
		})
	}
}
