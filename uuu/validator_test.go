package u_test

import (
	"encoding/json"
	"fmt"
	"testing"

	u "github.com/cachesdev/souuup/uuu"
)

// Define nested types used in tests
type NestedChild struct {
	Int int
}

type NestedStructure struct {
	Int    int
	String string
	Child  NestedChild
}

type Primitives struct {
	Float  float64
	Uint16 uint16
	Int32  int32
	Int    int
	Nested NestedStructure
}

type State struct {
	ID   int
	Name string
}

type City struct {
	ID    int
	Name  string
	State State
}

type District struct {
	ID   int
	Name string
}

type Club struct {
	ID          int
	Name        string
	CityID      int
	DistrictID  int
	Image       *string
	Address     *string
	MapLocation *string
	Contact     *string
	District    District
	City        City
}

func TestValidatePrimitives(t *testing.T) {
	req := Primitives{
		Float:  12.0,
		Uint16: 3,
		Int32:  3,
		Int:    1,
		Nested: NestedStructure{
			Int:    3,
			String: "ABC",
			Child: NestedChild{
				Int: 3,
			},
		},
	}

	v := u.NewSouuup(u.Schema{
		"myField":  u.Field(req.Float, u.MinN(10.0)),
		"myField2": u.Field(req.Uint16, u.MinN(uint16(3))),
		"myField3": u.Field(req.Int32, u.MinN(int32(3))),
		"myField4": u.Field(req.Int, u.MaxN(1)),
		"nestedField1": u.Schema{
			"IntValue":    u.Field(req.Nested.Int, u.MinN(3)),
			"StringValue": u.Field(req.Nested.String, u.MinS(3)),
			"Child": u.Schema{
				"Value": u.Field(req.Nested.Child.Int, u.MinN(3)),
			},
		},
	})

	if err := v.Validate(); err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}

func TestValidateComplexStructure(t *testing.T) {
	imageUrl := "https://example.com/image.jpg"
	address := "123 Main Street"
	mapLocation := "40.7128,-74.0060"
	contact := "John Doe: 555-1234"

	club := Club{
		ID:          1,
		Name:        "Sports Club",
		CityID:      2,
		DistrictID:  3,
		Image:       &imageUrl,
		Address:     &address,
		MapLocation: &mapLocation,
		Contact:     &contact,
		District: District{
			ID:   3,
			Name: "Downtown",
		},
		City: City{
			ID:   2,
			Name: "Metropolis",
			State: State{
				ID:   4,
				Name: "North State",
			},
		},
	}

	v := u.NewSouuup(u.Schema{
		"id":         u.Field(club.ID, u.MinN(1)),
		"name":       u.Field(club.Name, u.MinS(3)),
		"cityId":     u.Field(club.CityID, u.MinN(1)),
		"districtId": u.Field(club.DistrictID, u.MinN(1)),
		"image":      u.Field(*club.Image, u.MinS(5)),
		"address":    u.Field(*club.Address, u.MinS(5)),
		"district": u.Schema{
			"id":   u.Field(club.District.ID, u.MinN(1)),
			"name": u.Field(club.District.Name, u.MinS(3)),
		},
		"city": u.Schema{
			"id":   u.Field(club.City.ID, u.MinN(1)),
			"name": u.Field(club.City.Name, u.MinS(3)),
			"state": u.Schema{
				"id":   u.Field(club.City.State.ID, u.MinN(1)),
				"name": u.Field(club.City.State.Name, u.MinS(3)),
			},
		},
	})

	if err := v.Validate(); err != nil {
		t.Errorf("Expected validation to pass, got error: %v", err)
	}
}

func TestValidationWithInvalidValues(t *testing.T) {
	emptyString := ""
	shortAddress := "123"

	club := Club{
		ID:         0,
		Name:       "S",
		CityID:     0,
		DistrictID: 0,
		Image:      &emptyString,
		Address:    &shortAddress,
		District: District{
			ID:   0,
			Name: "",
		},
		City: City{
			ID:   0,
			Name: "",
			State: State{
				ID:   0,
				Name: "",
			},
		},
	}

	defaultStr := ""
	imageToValidate := defaultStr
	if club.Image != nil {
		imageToValidate = *club.Image
	}

	addressToValidate := defaultStr
	if club.Address != nil {
		addressToValidate = *club.Address
	}

	v := u.NewSouuup(u.Schema{
		"id":         u.Field(club.ID, u.MinN(1)),
		"name":       u.Field(club.Name, u.MinS(3)),
		"cityId":     u.Field(club.CityID, u.MinN(1)),
		"districtId": u.Field(club.DistrictID, u.MinN(1)),
		"image":      u.Field(imageToValidate, u.MinS(5)),
		"address":    u.Field(addressToValidate, u.MinS(5)),
		"district": u.Schema{
			"id":   u.Field(club.District.ID, u.MinN(1)),
			"name": u.Field(club.District.Name, u.MinS(3)),
		},
		"city": u.Schema{
			"id":   u.Field(club.City.ID, u.MinN(1)),
			"name": u.Field(club.City.Name, u.MinS(3)),
			"state": u.Schema{
				"id":   u.Field(club.City.State.ID, u.MinN(1)),
				"name": u.Field(club.City.State.Name, u.MinS(3)),
			},
		},
	})

	err := v.Validate()
	if err == nil {
		t.Fatal("Expected validation to fail, but it passed")
	}
}

func TestValidationWithNilPointers(t *testing.T) {
	club := Club{
		ID:          1,
		Name:        "Sports Club",
		CityID:      2,
		DistrictID:  3,
		Image:       nil,
		Address:     nil,
		MapLocation: nil,
		Contact:     nil,
		District: District{
			ID:   3,
			Name: "Downtown",
		},
		City: City{
			ID:   2,
			Name: "Metropolis",
			State: State{
				ID:   4,
				Name: "North State",
			},
		},
	}

	defaultStr := ""

	imageToValidate := defaultStr
	if club.Image != nil {
		imageToValidate = *club.Image
	}

	addressToValidate := defaultStr
	if club.Address != nil {
		addressToValidate = *club.Address
	}

	v := u.NewSouuup(u.Schema{
		"id":         u.Field(club.ID, u.MinN(1)),
		"name":       u.Field(club.Name, u.MinS(3)),
		"cityId":     u.Field(club.CityID, u.MinN(1)),
		"districtId": u.Field(club.DistrictID, u.MinN(1)),
		"image":      u.Field(imageToValidate, u.MinS(5)),
		"address":    u.Field(addressToValidate, u.MinS(5)),
		"district": u.Schema{
			"id":   u.Field(club.District.ID, u.MinN(1)),
			"name": u.Field(club.District.Name, u.MinS(3)),
		},
		"city": u.Schema{
			"id":   u.Field(club.City.ID, u.MinN(1)),
			"name": u.Field(club.City.Name, u.MinS(3)),
			"state": u.Schema{
				"id":   u.Field(club.City.State.ID, u.MinN(1)),
				"name": u.Field(club.City.State.Name, u.MinS(3)),
			},
		},
	})

	err := v.Validate()
	if err == nil {
		t.Fatal("Expected validation to fail due to nil pointers, but it passed")
	}
}

func TestComplexValidationFailure(t *testing.T) {
	shortImage := "img"
	shortAddress := "123"
	emptyContact := ""

	club := Club{
		ID:          -5,   // Invalid: negative number
		Name:        "AB", // Invalid: too short
		CityID:      0,    // Invalid: should be > 0
		DistrictID:  -2,   // Invalid: negative number
		Image:       &shortImage,
		Address:     &shortAddress,
		MapLocation: nil,
		Contact:     &emptyContact,
		District: District{
			ID:   0,   // Invalid: should be > 0
			Name: "D", // Invalid: too short
		},
		City: City{
			ID:   0,   // Invalid: should be > 0
			Name: "C", // Invalid: too short
			State: State{
				ID:   0,  // Invalid: should be > 0
				Name: "", // Invalid: empty string
			},
		},
	}

	defaultStr := ""

	imageToValidate := defaultStr
	if club.Image != nil {
		imageToValidate = *club.Image
	}

	addressToValidate := defaultStr
	if club.Address != nil {
		addressToValidate = *club.Address
	}

	contactToValidate := defaultStr
	if club.Contact != nil {
		contactToValidate = *club.Contact
	}

	mapLocationToValidate := defaultStr
	if club.MapLocation != nil {
		mapLocationToValidate = *club.MapLocation
	}

	v := u.NewSouuup(u.Schema{
		"id":          u.Field(club.ID, u.MinN(1), u.MaxN(1000)),
		"name":        u.Field(club.Name, u.MinS(3), u.MaxS(50)),
		"cityId":      u.Field(club.CityID, u.MinN(1)),
		"districtId":  u.Field(club.DistrictID, u.MinN(1)),
		"image":       u.Field(imageToValidate, u.MinS(5), u.MaxS(200)),
		"address":     u.Field(addressToValidate, u.MinS(5), u.MaxS(200)),
		"contact":     u.Field(contactToValidate, u.MinS(5)),
		"mapLocation": u.Field(mapLocationToValidate, u.MinS(5)),
		"district": u.Schema{
			"id":   u.Field(club.District.ID, u.MinN(1)),
			"name": u.Field(club.District.Name, u.MinS(3), u.MaxS(100)),
		},
		"city": u.Schema{
			"id":   u.Field(club.City.ID, u.MinN(1)),
			"name": u.Field(club.City.Name, u.MinS(3), u.MaxS(50)),
			"state": u.Schema{
				"id":   u.Field(club.City.State.ID, u.MinN(1)),
				"name": u.Field(club.City.State.Name, u.MinS(3), u.MaxS(50)),
			},
		},
	})

	err := v.Validate()
	if err == nil {
		t.Fatal("Expected validation to fail with multiple errors, but it passed")
	}

	fmt.Println("Validation Errors:")
	var prettyJSON map[string]any
	json.Unmarshal([]byte(err.Error()), &prettyJSON)
	prettyOutput, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Println(string(prettyOutput))

	t.Fail()
}
