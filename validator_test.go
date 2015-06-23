package validator_test

import (
	"fmt"
	"testing"

	"github.com/liuhengloveyou/validator"
)

type ValidateExample struct {
	IdCard      string `validate:"idcard"`
	Id2         int64  `validate:"idcard"`
	Name        string `validate:"nonzero"`
	Description string
	Age         int    `validate:"min=18"`
	Email       string `validate:"nonzero,email"`
	Address     struct {
		Street string `validate:"nonzero"`
		City   string `validate:"nonzero"`
	}
}

func TestParseTags(t *testing.T) {

	ve := ValidateExample{
		Name:        "Joe Doe", // valid as it's nonzero
		Description: "",        // valid no validation tag exists
		Age:         19,        // invalid as age is less than required 18
	}

	ve.Email = "@not.a.valid.email"
	ve.Address.City = "Some City" // valid
	ve.Address.Street = ""        // invalid
	ve.IdCard = "41289719870908876x"
	ve.Id2 = 123

	err := validator.Validate(ve)
	fmt.Println(err)
}

func TestRegion(t *testing.T) {
	type Demo struct {
		Region string `validate:"region"`
	}

	ve := Demo{
		Region: "412826"}

	err := validator.Validate(ve)
	fmt.Println(err)

	ve = Demo{
		Region: "110000"}

	err = validator.Validate(ve)
	fmt.Println(err)
}

func TestCN(t *testing.T) {
	type Demo struct {
		Name string `validate:"unicn"`
		Age  int32  `validate:"unicn"`
	}

	v1 := &Demo{
		Name: "小明",
		Age:  11,
	}

	e := validator.Validate(v1)
	fmt.Println(e)
}

func TestNoneor(t *testing.T) {
	type Demo struct {
		Name string `validate:"unicn"`
		Age  int32  `validate:"noneor"`
		// Age  int32  `validate:"noneor,min=18"`
	}

	v1 := &Demo{
		Name: "小明",
		Age:  11,
	}

	e := validator.Validate(v1)
	fmt.Println(e)
}
