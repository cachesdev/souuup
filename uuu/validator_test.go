package u_test

import (
	"fmt"
	"testing"

	u "github.com/cachesdev/souuup/uuu"
)

func TestValidate(t *testing.T) {
	type C struct {
		D int
	}

	type NestedPair struct {
		A int
		B string
		C C
	}

	req := struct {
		NumF       float64
		NumUint    uint16
		Num32      int32
		Num        int
		NestedPair NestedPair
	}{
		NumF:    3.0,
		NumUint: 3,
		Num32:   3,
		Num:     3,
		NestedPair: NestedPair{
			A: 1,
			B: "A",
			C: C{
				D: 1,
			},
		},
	}

	v := u.Souuup{
		"myField":  u.Field(req.NumF, u.MinN(10.0)),
		"myField2": u.Field(req.NumUint, u.MinN(uint16(3))),
		"myField3": u.Field(req.Num32, u.MinN(int32(3))),
		"myField4": u.Field(req.Num, u.MaxN(1)),
		"nestedField1": u.Nested(u.Souuup{
			"A": u.Field(req.NestedPair.A, u.MinN(3)),
			"B": u.Field(req.NestedPair.B, u.MinS(3)),
			"C": u.Nested(u.Souuup{
				"D": u.Field(req.NestedPair.C.D, u.MinN(3)),
			}),
		}),
	}

	if err := v.ValidateSouuup(); err != nil {
		fmt.Println(err.Error())
		t.Error("Expected validation to pass")
	}
}

type State struct {
	ID               int
	StateDescription string
}

type City struct {
	ID              int
	CityDescription string
	State           State
}

type District struct {
	ID                  int
	DistrictDescription string
}

type Club struct {
	ID              int
	ClubDescription string
	CityID          int
	DistrictID      int
	Image           *string
	Address         *string
	MapLocation     *string
	Contact         *string
	District        District
	City            City
}

func TestValidateComplex(t *testing.T) {
	imageUrl := "https://example.com/image.jpg"
	address := "123 Main Street"
	mapLocation := "40.7128,-74.0060"
	contact := "John Doe: 555-1234"

	club := Club{
		ID:              1,
		ClubDescription: "Sports Club",
		CityID:          2,
		DistrictID:      3,
		Image:           &imageUrl,
		Address:         &address,
		MapLocation:     &mapLocation,
		Contact:         &contact,
		District: District{
			ID:                  3,
			DistrictDescription: "Downtown",
		},
		City: City{
			ID:              2,
			CityDescription: "Metropolis",
			State: State{
				ID:               4,
				StateDescription: "North State",
			},
		},
	}

	v := u.Souuup{
		"id":          u.Field(club.ID, u.MinN(1)),
		"description": u.Field(club.ClubDescription, u.MinS(3)),
		"cityId":      u.Field(club.CityID, u.MinN(1)),
		"districtId":  u.Field(club.DistrictID, u.MinN(1)),
		"image":       u.Field(*club.Image, u.MinS(5)),
		"address":     u.Field(*club.Address, u.MinS(5)),
		"district": u.Nested(u.Souuup{
			"id":          u.Field(club.District.ID, u.MinN(1)),
			"description": u.Field(club.District.DistrictDescription, u.MinS(3)),
		}),
		"city": u.Nested(u.Souuup{
			"id":          u.Field(club.City.ID, u.MinN(1)),
			"description": u.Field(club.City.CityDescription, u.MinS(3)),
			"state": u.Nested(u.Souuup{
				"id":          u.Field(club.City.State.ID, u.MinN(1)),
				"description": u.Field(club.City.State.StateDescription, u.MinS(3)),
			}),
		}),
	}

	if err := v.ValidateSouuup(); err != nil {
		fmt.Println(err.Error())
		t.Error("Expected validation to pass")
	}
}

func TestValidationFailure(t *testing.T) {
	// Create a club instance with INVALID data to trigger failures
	emptyString := ""
	shortAddress := "123" // Too short, will fail MinS(5)

	club := Club{
		ID:              0,             // Will fail MinN(1)
		ClubDescription: "S",           // Too short, will fail MinS(3)
		CityID:          0,             // Will fail MinN(1)
		DistrictID:      0,             // Will fail MinN(1)
		Image:           &emptyString,  // Empty, will fail MinS(5)
		Address:         &shortAddress, // Too short, will fail MinS(5)
		MapLocation:     nil,           // Nil pointer, will produce a separate error
		Contact:         nil,           // Nil pointer, will produce a separate error
		District: District{
			ID:                  0,  // Will fail MinN(1)
			DistrictDescription: "", // Empty, will fail MinS(3)
		},
		City: City{
			ID:              0,  // Will fail MinN(1)
			CityDescription: "", // Empty, will fail MinS(3)
			State: State{
				ID:               0,  // Will fail MinN(1)
				StateDescription: "", // Empty, will fail MinS(3)
			},
		},
	}

	// Validate the club, expecting multiple failures
	v := u.Souuup{
		"id":          u.Field(club.ID, u.MinN(1)),
		"description": u.Field(club.ClubDescription, u.MinS(3)),
		"cityId":      u.Field(club.CityID, u.MinN(1)),
		"districtId":  u.Field(club.DistrictID, u.MinN(1)),
		"image":       u.Field(*club.Image, u.MinS(5)),
		"address":     u.Field(*club.Address, u.MinS(5)),
		"district": u.Nested(u.Souuup{
			"id":          u.Field(club.District.ID, u.MinN(1)),
			"description": u.Field(club.District.DistrictDescription, u.MinS(3)),
		}),
		"city": u.Nested(u.Souuup{
			"id":          u.Field(club.City.ID, u.MinN(1)),
			"description": u.Field(club.City.CityDescription, u.MinS(3)),
			"state": u.Nested(u.Souuup{
				"id":          u.Field(club.City.State.ID, u.MinN(1)),
				"description": u.Field(club.City.State.StateDescription, u.MinS(3)),
			}),
		}),
	}

	err := v.ValidateSouuup()
	if err == nil {
		t.Fatal("Expected validation to fail, but it passed")
	}

	// Print the error to see the JSON output format
	fmt.Println("Validation Error JSON:")
	fmt.Println(err.Error())
}
