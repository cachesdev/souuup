// main package of Soup
package u

import (
	"golang.org/x/exp/constraints"
)

// FieldTag represents the "key" of a field, and will be used to identify a field on
// an error map
type FieldTag = string

// a Souuup instance
type Souuup map[FieldTag]Validable

// a Souuup instance for boring people
type Validator = Souuup

type Numeric interface {
	constraints.Signed | constraints.Float
}

// Valid runs validation on every field, and sets the error map.
func (v Souuup) Valid() bool {
	for tag, field := range v {
		field.SetTag(tag)
		if !field.Validate() {
			return false
		}
	}
	return true
}

type Validable interface {
	Validate() bool
	Tag() string
	SetTag(string)
}
