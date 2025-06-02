package r_test

import (
	"testing"

	"github.com/cachesdev/souuup/internal/testutil"
	"github.com/cachesdev/souuup/r"
	"github.com/cachesdev/souuup/u"
)

func TestMinN(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		min      int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "equal to minimum",
			value:   18,
			min:     18,
			wantErr: false,
		},
		{
			name:    "greater than minimum",
			value:   25,
			min:     18,
			wantErr: false,
		},
		{
			name:     "less than minimum",
			value:    15,
			min:      18,
			wantErr:  true,
			errorMsg: "value is 15, but needs to be at least 18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[int]{Value: tt.value}
			err := r.MinN(tt.min)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}

	// Test with float type
	t.Run("float: greater than minimum", func(t *testing.T) {
		fs := u.FieldState[float64]{Value: 3.14}
		err := r.MinN(3.0)(fs)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("float: less than minimum", func(t *testing.T) {
		fs := u.FieldState[float64]{Value: 2.5}
		err := r.MinN(3.0)(fs)
		expected := "value is 2.5, but needs to be at least 3"
		if err == nil {
			t.Errorf("expected error but got nil")
		} else if err.Error() != expected {
			t.Errorf("expected error message: %q, got: %q", expected, err.Error())
		}
	})
}

func TestMaxN(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		max      int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "equal to maximum",
			value:   100,
			max:     100,
			wantErr: false,
		},
		{
			name:    "less than maximum",
			value:   80,
			max:     100,
			wantErr: false,
		},
		{
			name:     "greater than maximum",
			value:    120,
			max:      100,
			wantErr:  true,
			errorMsg: "value is 120, but needs to be at most 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[int]{Value: tt.value}
			err := r.MaxN(tt.max)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestGt(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		limit    int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "greater than limit",
			value:   19,
			limit:   18,
			wantErr: false,
		},
		{
			name:     "equal to limit",
			value:    18,
			limit:    18,
			wantErr:  true,
			errorMsg: "value is 18, but needs to be greater than 18",
		},
		{
			name:     "less than limit",
			value:    17,
			limit:    18,
			wantErr:  true,
			errorMsg: "value is 17, but needs to be greater than 18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[int]{Value: tt.value}
			err := r.Gt(tt.limit)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestGte(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		limit    int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "greater than limit",
			value:   19,
			limit:   18,
			wantErr: false,
		},
		{
			name:    "equal to limit",
			value:   18,
			limit:   18,
			wantErr: false,
		},
		{
			name:     "less than limit",
			value:    17,
			limit:    18,
			wantErr:  true,
			errorMsg: "value is 17, but needs to be at least 18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[int]{Value: tt.value}
			err := r.Gte(tt.limit)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestLt(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		limit    int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "less than limit",
			value:   99,
			limit:   100,
			wantErr: false,
		},
		{
			name:     "equal to limit",
			value:    100,
			limit:    100,
			wantErr:  true,
			errorMsg: "value is 100, but needs to be less than 100",
		},
		{
			name:     "greater than limit",
			value:    101,
			limit:    100,
			wantErr:  true,
			errorMsg: "value is 101, but needs to be less than 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[int]{Value: tt.value}
			err := r.Lt(tt.limit)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestLte(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		limit    int
		wantErr  bool
		errorMsg string
	}{
		{
			name:    "less than limit",
			value:   99,
			limit:   100,
			wantErr: false,
		},
		{
			name:    "equal to limit",
			value:   100,
			limit:   100,
			wantErr: false,
		},
		{
			name:     "greater than limit",
			value:    101,
			limit:    100,
			wantErr:  true,
			errorMsg: "value is 101, but needs to be at most 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[int]{Value: tt.value}
			err := r.Lte(tt.limit)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}

func TestNeqN(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		notEqual int
		wantErr  bool
		errorMsg string
	}{
		{
			name:     "equal value",
			value:    42,
			notEqual: 42,
			wantErr:  true,
			errorMsg: "value is 42, but needs to not equal to 42",
		},
		{
			name:     "different value",
			value:    42,
			notEqual: 24,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := u.FieldState[int]{Value: tt.value}
			err := r.NeqN(tt.notEqual)(fs)
			testutil.CheckError(t, err, tt.wantErr, tt.errorMsg)
		})
	}
}
