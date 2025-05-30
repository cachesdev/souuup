package r

import (
	"fmt"

	"github.com/cachesdev/souuup/u"
)

// MinS validates if a string's length is at least n characters.
//
// Example:
//
//	// Validate that a name is at least 2 characters long
//	nameField := u.Field("John", u.MinS(2))
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
//	usernameField := u.Field("john doe", u.MaxS(20))
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
//	otpField := u.Field("123456", u.LenS(6))
func LenS(n int) u.StringRule {
	return func(fd u.FieldState[string]) error {
		if len(fd.Value) != n {
			return fmt.Errorf("length is %d, but needs to be exactly %d", len(fd.Value), n)
		}
		return nil
	}
}
