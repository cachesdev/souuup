package r_test

import (
	"testing"

	"github.com/cachesdev/souuup/internal/testutil"
	"github.com/cachesdev/souuup/r"
	"github.com/cachesdev/souuup/u"
)

func TestMinS(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		min      int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "exact minimum length",
			value:   "ab",
			min:     2,
			wantErr: false,
		},
		{
			name:    "greater than minimum length",
			value:   "abcde",
			min:     3,
			wantErr: false,
		},
		{
			name:     "less than minimum length",
			value:    "a",
			min:      2,
			wantErr:  true,
			errorMsg: "length is 1, but needs to be at least 2",
		},
		{
			name:     "empty string",
			value:    "",
			min:      1,
			wantErr:  true,
			errorMsg: "length is 0, but needs to be at least 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[string]{Value: tt.value}
			err := r.MinS(tt.min)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestMaxS(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		max      int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "exact maximum length",
			value:   "abcde",
			max:     5,
			wantErr: false,
		},
		{
			name:    "less than maximum length",
			value:   "abc",
			max:     5,
			wantErr: false,
		},
		{
			name:     "greater than maximum length",
			value:    "abcdef",
			max:      5,
			wantErr:  true,
			errorMsg: "length is 6, but needs to be at most 5",
		},
		{
			name:    "empty string",
			value:   "",
			max:     5,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[string]{Value: tt.value}
			err := r.MaxS(tt.max)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestLenS(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		length   int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "exact length",
			value:   "abcde",
			length:  5,
			wantErr: false,
		},
		{
			name:     "shorter than required",
			value:    "abc",
			length:   5,
			wantErr:  true,
			errorMsg: "length is 3, but needs to be exactly 5",
		},
		{
			name:     "longer than required",
			value:    "abcdefg",
			length:   5,
			wantErr:  true,
			errorMsg: "length is 7, but needs to be exactly 5",
		},
		{
			name:     "empty string",
			value:    "",
			length:   5,
			wantErr:  true,
			errorMsg: "length is 0, but needs to be exactly 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[string]{Value: tt.value}
			err := r.LenS(tt.length)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestInS(t *testing.T) {
	validValues := []string{"small", "medium", "large"}

	tests := []struct {
		name     string
		value    string
		set      []string
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "value in set",
			value:   "medium",
			set:     validValues,
			wantErr: false,
		},
		{
			name:     "value not in set",
			value:    "extra-large",
			set:      validValues,
			wantErr:  true,
			errorMsg: `"extra-large" is not in [small medium large], but should be`,
		},
		{
			name:     "empty value with non-empty set",
			value:    "",
			set:      validValues,
			wantErr:  true,
			errorMsg: `"" is not in [small medium large], but should be`,
		},
		{
			name:     "empty set (should not be used in practice)",
			value:    "anything",
			set:      []string{},
			wantErr:  true,
			errorMsg: `"anything" is not in [], but should be`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[string]{Value: tt.value}
			err := r.InS(tt.set)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestNotInS(t *testing.T) {
	invalidValues := []string{"rejected", "invalid", "error"}

	tests := []struct {
		name     string
		value    string
		set      []string
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "value not in set",
			value:   "approved",
			set:     invalidValues,
			wantErr: false,
		},
		{
			name:     "value in set",
			value:    "rejected",
			set:      invalidValues,
			wantErr:  true,
			errorMsg: `"rejected" is in [rejected invalid error], but shouldn't be`,
		},
		{
			name:    "empty value with set not containing empty string",
			value:   "",
			set:     invalidValues,
			wantErr: false,
		},
		{
			name:    "empty set (should not be used in practice)",
			value:   "anything",
			set:     []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[string]{Value: tt.value}
			err := r.NotInS(tt.set)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestContainsS(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		substr   string
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "contains substring",
			value:   "123 London Street",
			substr:  "Street",
			wantErr: false,
		},
		{
			name:     "does not contain substring",
			value:    "123 London Road",
			substr:   "Street",
			wantErr:  true,
			errorMsg: `"123 London Road" does not contain "Street", but needs to`,
		},
		{
			name:    "substring is empty (always contained)",
			value:   "Any string",
			substr:  "",
			wantErr: false,
		},
		{
			name:     "value is empty",
			value:    "",
			substr:   "something",
			wantErr:  true,
			errorMsg: `"" does not contain "something", but needs to`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[string]{Value: tt.value}
			err := r.ContainsS(tt.substr)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}
