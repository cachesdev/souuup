package r

import (
	"fmt"
	"slices"
	"strings"

	"github.com/cachesdev/souuup/u"
)

// Every validates that every element in the slice satisfies the given rule.
// It returns an error as soon as any element fails the validation.
//
// Example:
//
//	// Validate that all interests are at least 3 characters long
//	interestsField := u.Field(user.Interests, r.Every(r.MinS(3)))
func Every[T any](rule u.Rule[T]) u.SliceRule[T] {
	return func(fs u.FieldState[[]T]) error {
		slice := fs.Value
		if len(slice) == 0 {
			return nil // Empty slices pass validation by default
		}

		errors := make(map[int]error)
		for i, item := range slice {
			itemState := u.FieldState[T]{Value: item}
			if err := rule(itemState); err != nil {
				errors[i] = err
			}
		}

		if len(errors) > 0 {
			var errMsgs []string
			errMsgs = append(errMsgs, fmt.Sprintf("%d elements failed validation", len(errors)))
			for i, err := range errors {
				errMsgs = append(errMsgs, fmt.Sprintf("  [%d]: %s", i, err.Error()))
			}
			return fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
		}
		return nil
	}
}

// Some validates that at least one element in the slice satisfies the given rule.
// It returns an error if all elements fail the validation.
//
// Example:
//
//	// Validate that at least one admin exists
//	adminsField := u.Field(users, r.Some(func(fs u.FieldState[User]) error {
//		if !fs.Value.IsAdmin {
//			return fmt.Errorf("user is not an admin")
//		}
//		return nil
//	}))
func Some[T any](rule u.Rule[T]) u.SliceRule[T] {
	return func(fs u.FieldState[[]T]) error {
		slice := fs.Value
		if len(slice) == 0 {
			return fmt.Errorf("slice is empty, but needs at least one valid element")
		}

		somePassed := false

		errors := make(map[int]error)
		for i, item := range slice {
			itemState := u.FieldState[T]{Value: item}
			if err := rule(itemState); err != nil {
				errors[i] = err
			} else {
				somePassed = true
			}
		}

		if !somePassed {
			var errMsgs []string
			errMsgs = append(errMsgs, fmt.Sprintf("all %d elements failed validation", len(errors)))
			for i, err := range errors {
				errMsgs = append(errMsgs, fmt.Sprintf("  [%d]: %s", i, err.Error()))
			}
			return fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
		}

		return nil
	}
}

// None validates that no element in the slice satisfies the given rule.
// It returns an error if any element passes the validation.
//
// Example:
//
//	// Validate that no banned words are used
//	wordsField := u.Field(words, r.None(func(fs u.FieldState[string]) error {
//		bannedWords := map[string]bool{"forbidden": true, "banned": true}
//		if bannedWords[fs.Value] {
//			return nil // If word is banned, the inner rule passes (no error)
//		}
//		return fmt.Errorf("word is not banned") // If word is not banned, inner rule fails
//	}))
func None[T any](rule u.Rule[T]) u.SliceRule[T] {
	return func(fs u.FieldState[[]T]) error {
		slice := fs.Value
		if len(slice) == 0 {
			return nil // Empty slices pass validation by default
		}

		unexpectedPassIndices := []int{}

		for i, item := range slice {
			itemState := u.FieldState[T]{Value: item}
			if rule(itemState) == nil {
				unexpectedPassIndices = append(unexpectedPassIndices, i)
			}
		}

		if len(unexpectedPassIndices) > 0 {
			var errMsgs []string
			errMsgs = append(errMsgs, fmt.Sprintf("%d elements unexpectedly passed validation (expected none to pass rule)", len(unexpectedPassIndices)))

			for _, i := range unexpectedPassIndices {
				errMsgs = append(errMsgs, fmt.Sprintf("  [%d]: failed", i))
			}

			return fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
		}
		return nil
	}
}

// MinLen validates that a slice has at least n elements.
//
// Example:
//
//	// Validate that a user has at least one interest
//	interestsField := u.Field(user.Interests, r.MinLen(1))
func MinLen[T any](n int) u.SliceRule[T] {
	return func(fs u.FieldState[[]T]) error {
		length := len(fs.Value)
		if length < n {
			return fmt.Errorf("length is %d, but needs to be at least %d", length, n)
		}
		return nil
	}
}

// MaxLen validates that a slice has at most n elements.
//
// Example:
//
//	// Validate that a user has at most 5 interests
//	interestsField := u.Field(user.Interests, r.MaxLen(5))
func MaxLen[T any](n int) u.SliceRule[T] {
	return func(fs u.FieldState[[]T]) error {
		length := len(fs.Value)
		if length > n {
			return fmt.Errorf("length is %d, but needs to be at most %d", length, n)
		}
		return nil
	}
}

// ExactLen validates that a slice has exactly n elements.
//
// Example:
//
//	// Validate that a team has exactly 5 members
//	teamMembersField := u.Field(team.Members, r.ExactLen(5))
func ExactLen[T any](n int) u.SliceRule[T] {
	return func(fs u.FieldState[[]T]) error {
		length := len(fs.Value)
		if length != n {
			return fmt.Errorf("length is %d, but needs to be exactly %d", length, n)
		}
		return nil
	}
}

// Contains validates that a slice contains a matching element for comparable slices.
//
// Example:
//
//	// Validate that a team has a goalkeeper
//	teamMembersField := u.Field(team.Members, r.Contains("GK"))
func Contains[T comparable](member T) u.SliceRule[T] {
	return func(fs u.FieldState[[]T]) error {
		if !slices.Contains(fs.Value, member) {
			return fmt.Errorf("%v does not contain %v, but needs to", fs.Value, member)
		}
		return nil
	}
}
