package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func Validate(data any) map[string]string {
	err := validate.Struct(data)
	if err != nil {
		errMaps := make(map[string]string)

		if validationErr, ok := err.(validator.ValidationErrors); ok {
			for _, vErr := range validationErr {
				errMaps[vErr.Field()] = fmt.Sprintf("failed on '%s' rule", vErr.Tag())
			}
		} else {
			errMaps["error"] = err.Error()
		}
		return errMaps
	}
	return nil
}