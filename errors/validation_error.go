package errors

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorResponse(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]string)
		for _, fe := range ve {
			out[fe.Field()] = ValidationErrorMessage(fe)
		}
		data, err := json.Marshal(out)
		if err != nil {
			return "invalid request"
		}
		return string(data)

	}
	return "invalid request"
}

func ValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	}
	return "is invalid"
}
