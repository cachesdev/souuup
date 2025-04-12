package u_test

import (
	"testing"

	u "github.com/cachesdev/souuup/uuu"
)

func TestValidate(t *testing.T) {
	req := struct {
		NumF    float64
		NumUint uint16
		Num32   int32
		Num     int
	}{
		NumF:    3.0,
		NumUint: 3,
		Num32:   3,
		Num:     3,
	}

	v := u.Souuup{
		"myField":  u.Field(req.NumF, u.MinN(3.0)),
		"myField2": u.Field(req.NumUint, u.MinN(uint16(3))),
		"myField3": u.Field(req.Num32, u.MinN(int32(3))),
		"myField4": u.Field(req.Num, u.MinN(3)),
	}

	if !v.Valid() {
		t.Error("Expected validation to pass")
	}
}
