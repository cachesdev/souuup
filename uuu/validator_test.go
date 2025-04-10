package u_test

import (
	"testing"

	u "github.com/cachesdev/souuup/uuu"
)

func TestValidate(t *testing.T) {
	v := u.Souuup{
		"myField": u.Field("test", u.MinS(5)),
	}

	if !v.Valid() {
		t.Error("Expected validation to pass")
	}
}
