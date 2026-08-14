package validator

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func Struct(s any) error {
	return Validate.Struct(s)
}
