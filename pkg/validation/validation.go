package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(data any) error {
	err := validate.Struct(data)
	if err != nil {
		if validationErros, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErros {
				return fmt.Errorf("%s %s", fieldErr.Field(), fieldErr.Error())
			}
		}

	}
	return nil
}