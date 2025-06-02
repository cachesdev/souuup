package r

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cachesdev/souuup/u"
)

// MinS validates if a string's length is at least n characters.
//
// Example:
//
//	// Validate that a name is at least 2 characters long
//	nameField := u.Field("John", r.MinS(2))
func MinS(n int) u.StringRule {
	return func(fd u.FieldState[string]) error {
		if len(fd.Value) < n {
			return fmt.Errorf("length is %d, but needs to be at least %d", len(fd.Value), n)
		}
		return nil
	}
}

// MaxS validates if a string's length is at most n characters.
//
// Example:
//
//	// Validate that a username is at most 20 characters long
//	usernameField := u.Field("john doe", r.MaxS(20))
func MaxS(n int) u.StringRule {
	return func(fd u.FieldState[string]) error {
		if len(fd.Value) > n {
			return fmt.Errorf("length is %d, but needs to be at most %d", len(fd.Value), n)
		}
		return nil
	}
}

// LenS validates if a string's length is exactly n characters.
//
// Example:
//
//	// Validate that a passcode is exactly 6 characters long
//	otpField := u.Field("123456", r.LenS(6))
func LenS(n int) u.StringRule {
	return func(fd u.FieldState[string]) error {
		if len(fd.Value) != n {
			return fmt.Errorf("length is %d, but needs to be exactly %d", len(fd.Value), n)
		}
		return nil
	}
}

// InS validates if a string is contained within a set of strings
//
// Example:
//
//	// Validate that a size is small, medium, or large
//	sizeField := u.Field("small", r.InS(["small", "medium", "large"]))
func InS(set []string) u.StringRule {
	return func(fs u.FieldState[string]) error {
		if !slices.Contains(set, fs.Value) {
			return fmt.Errorf("%q is not in %v, but should be", fs.Value, set)
		}
		return nil
	}
}

// NotInS validates if a string is not contained within a set of strings
//
// Example:
//
//	// Validate that a status is not rejected or invalid
//	statusField := u.Field("completed", r.NotInS(["rejected", "invalid"]))
func NotInS(set []string) u.StringRule {
	return func(fs u.FieldState[string]) error {
		if slices.Contains(set, fs.Value) {
			return fmt.Errorf("%q is in %v, but shouldn't be", fs.Value, set)
		}
		return nil
	}
}

// ContainsS validates if a string contains a substring
//
// Example:
//
//	// Validate that an address contains "Street"
//	addrField := u.Field("123 London Street", r.ContainsS("Street"))
func ContainsS(substr string) u.StringRule {
	return func(fs u.FieldState[string]) error {
		if !strings.Contains(fs.Value, substr) {
			return fmt.Errorf("%q does not contain %q, but needs to", fs.Value, substr)
		}
		return nil
	}
}
