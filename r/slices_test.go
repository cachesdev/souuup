package r_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/cachesdev/souuup/r"
	"github.com/cachesdev/souuup/u"
)

func TestMinLen(t *testing.T) {
	tests := []struct {
		name     string
		value    []string
		min      int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "exact minimum length",
			value:   []string{"item1", "item2"},
			min:     2,
			wantErr: false,
		},
		{
			name:    "greater than minimum length",
			value:   []string{"item1", "item2", "item3"},
			min:     2,
			wantErr: false,
		},
		{
			name:     "less than minimum length",
			value:    []string{"item1"},
			min:      2,
			wantErr:  true,
			errorMsg: "length is 1, but needs to be at least 2",
		},
		{
			name:     "empty slice",
			value:    []string{},
			min:      1,
			wantErr:  true,
			errorMsg: "length is 0, but needs to be at least 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[[]string]{Value: tt.value}
			err := r.MinLen[string](tt.min)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestMaxLen(t *testing.T) {
	tests := []struct {
		name     string
		value    []int
		max      int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "exact maximum length",
			value:   []int{1, 2, 3, 4, 5},
			max:     5,
			wantErr: false,
		},
		{
			name:    "less than maximum length",
			value:   []int{1, 2, 3},
			max:     5,
			wantErr: false,
		},
		{
			name:     "greater than maximum length",
			value:    []int{1, 2, 3, 4, 5, 6},
			max:      5,
			wantErr:  true,
			errorMsg: "length is 6, but needs to be at most 5",
		},
		{
			name:    "empty slice",
			value:   []int{},
			max:     5,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[[]int]{Value: tt.value}
			err := r.MaxLen[int](tt.max)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestExactLen(t *testing.T) {
	tests := []struct {
		name     string
		value    []float64
		length   int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "exact length",
			value:   []float64{1.1, 2.2, 3.3, 4.4, 5.5},
			length:  5,
			wantErr: false,
		},
		{
			name:     "shorter than required",
			value:    []float64{1.1, 2.2, 3.3},
			length:   5,
			wantErr:  true,
			errorMsg: "length is 3, but needs to be exactly 5",
		},
		{
			name:     "longer than required",
			value:    []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6, 7.7},
			length:   5,
			wantErr:  true,
			errorMsg: "length is 7, but needs to be exactly 5",
		},
		{
			name:     "empty slice",
			value:    []float64{},
			length:   5,
			wantErr:  true,
			errorMsg: "length is 0, but needs to be exactly 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[[]float64]{Value: tt.value}
			err := r.ExactLen[float64](tt.length)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		value    []string
		member   string
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "contains member",
			value:   []string{"apple", "banana", "cherry"},
			member:  "banana",
			wantErr: false,
		},
		{
			name:     "does not contain member",
			value:    []string{"apple", "cherry", "date"},
			member:   "banana",
			wantErr:  true,
			errorMsg: "[apple cherry date] does not contain banana, but needs to",
		},
		{
			name:     "empty slice",
			value:    []string{},
			member:   "anything",
			wantErr:  true,
			errorMsg: "[] does not contain anything, but needs to",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[[]string]{Value: tt.value}
			err := r.Contains(tt.member)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestEvery(t *testing.T) {
	minLength3 := func(fs u.FieldState[string]) error {
		if len(fs.Value) < 3 {
			return fmt.Errorf("length is %d, but needs to be at least 3", len(fs.Value))
		}
		return nil
	}

	tests := []struct {
		name     string
		value    []string
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "all elements pass",
			value:   []string{"abc", "defg", "hijkl"},
			wantErr: false,
		},
		{
			name:    "empty slice passes",
			value:   []string{},
			wantErr: false,
		},
		{
			name:     "one element fails",
			value:    []string{"abc", "de", "fghi"},
			wantErr:  true,
			errorMsg: "1 elements failed validation",
		},
		{
			name:     "multiple elements fail",
			value:    []string{"ab", "cd", "ef"},
			wantErr:  true,
			errorMsg: "3 elements failed validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[[]string]{Value: tt.value}
			err := r.Every(minLength3)(fs)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err != nil)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message containing: %q, got: %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestSome(t *testing.T) {
	isEven := func(fs u.FieldState[int]) error {
		if fs.Value%2 != 0 {
			return fmt.Errorf("%d is not even", fs.Value)
		}
		return nil
	}

	tests := []struct {
		name     string
		value    []int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "at least one element passes",
			value:   []int{1, 3, 5, 6, 7},
			wantErr: false,
		},
		{
			name:    "all elements pass",
			value:   []int{2, 4, 6, 8},
			wantErr: false,
		},
		{
			name:     "no elements pass",
			value:    []int{1, 3, 5, 7, 9},
			wantErr:  true,
			errorMsg: "all 5 elements failed validation",
		},
		{
			name:     "empty slice",
			value:    []int{},
			wantErr:  true,
			errorMsg: "slice is empty, but needs at least one valid element",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[[]int]{Value: tt.value}
			err := r.Some(isEven)(fs)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err != nil)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message containing: %q, got: %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestNone(t *testing.T) {
	isNegative := func(fs u.FieldState[int]) error {
		if fs.Value < 0 {
			return nil // The rule passes for negative numbers
		}
		return fmt.Errorf("%d is not negative", fs.Value)
	}

	tests := []struct {
		name     string
		value    []int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "no element passes (all positive)",
			value:   []int{1, 3, 5, 7},
			wantErr: false,
		},
		{
			name:     "one element passes (has negative)",
			value:    []int{1, -3, 5, 7},
			wantErr:  true,
			errorMsg: "1 elements unexpectedly passed validation (expected none to pass rule)",
		},
		{
			name:     "multiple elements pass (has negatives)",
			value:    []int{-1, 3, -5, 7},
			wantErr:  true,
			errorMsg: "2 elements unexpectedly passed validation (expected none to pass rule)",
		},
		{
			name:    "empty slice",
			value:   []int{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[[]int]{Value: tt.value}
			err := r.None(isNegative)(fs)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error: %v, got: %v", tt.wantErr, err != nil)
				return
			}

			if tt.wantErr && err != nil {
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message containing: %q, got: %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}
