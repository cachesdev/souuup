package testutil

import (
	"testing"
)

// CheckError is a helper function to check if an error matches expected.
// It validates both the error presence and, when appropriate, the error message.
func CheckError(t *testing.T, err error, wantErr bool, errorMsg string) {
	t.Helper()

	if (err != nil) != wantErr {
		t.Errorf("expected error: %v, got: %v", wantErr, err != nil)
		return
	}

	if wantErr && err != nil && err.Error() != errorMsg {
		t.Errorf("expected error message: %q, got: %q", errorMsg, err.Error())
	}
}
