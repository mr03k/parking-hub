package validators

import (
	"github.com/go-playground/validator/v10"
	"strconv"
)

type TimestampValidator struct {
}

func NewTimestampValidator() TimestampValidator {
	return TimestampValidator{}
}

// timestamp validator
func (uv TimestampValidator) Handler() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().Int()
		return len(strconv.Itoa(int(value))) == 10
	}
}
