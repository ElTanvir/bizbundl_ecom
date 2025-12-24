package util

import (
	"fmt"
	"reflect"
	"strings"
	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func RegisterTagName() {
	Validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func ValidateStruct(s any) (map[string]string, error) {
	err := Validator.Struct(s)
	if err == nil {
		return nil, nil
	}

	errors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		message := getErrorMsg(err)
		errors[field] = message
	}

	return errors, err
}

func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Should be at least %s characters long", err.Param())
	default:
		return fmt.Sprintf("Invalid value: %s", err.Tag())
	}
}
