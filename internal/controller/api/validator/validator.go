package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

const partsInJSONTag = 2

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func NewValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		if name := fld.Tag.Get("param"); name != "" {
			return name
		}

		if name := fld.Tag.Get("query"); name != "" {
			return name
		}

		if jsonTag := fld.Tag.Get("json"); jsonTag != "" {
			name := strings.SplitN(jsonTag, ",", partsInJSONTag)[0]
			return name
		}

		return fld.Name
	})

	return v
}
