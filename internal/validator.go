package internal

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var (
	Validator = validator.New()
)

func Validate(s interface{}) error {
	// CHeck for errors
	err := Validator.Struct(s)

	// Check if there is any validation error
	if err == nil {
		return nil
	}

	// Map the first validation error as ErrInvalid and return it
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return &Error{
			Code:    ErrInternal,
			Message: "cannot convert error to ValidationErrors",
		}
	}
	// Get first validation error
	validationError := validationErrors[0]

	// Map validation error to ErrInvalid
	switch validationError.Tag() {
	case "required":
		return &Error{
			Code:    ErrInvalid,
			Message: fmt.Sprintf("the field %s is required", validationError.Field()),
			Err:     err,
		}
	case "url":
		return &Error{
			Code:    ErrInvalid,
			Message: fmt.Sprintf("the field %s must be a valid URL", validationError.Field()),
			Err:     err,
		}
	default:
		return &Error{Code: ErrInvalid}
	}
}
