package utils

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

// Function to translate validation errors to human-readable messages using JSON field names
func formatValidationError(err error, obj interface{}) map[string]string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": "Invalid input data"}
	}

	errors := make(map[string]string)

	// Get the reflect type of the struct
	val := reflect.TypeOf(obj)

	for _, fieldError := range validationErrors {
		field, _ := val.Elem().FieldByName(fieldError.Field())

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = fieldError.Field()
		}

		// Customize error messages based on validation tag
		switch fieldError.Tag() {
		case "required":
			errors[jsonTag] = fmt.Sprintf("%s is required", jsonTag)
		case "email":
			errors[jsonTag] = fmt.Sprintf("%s must be a valid email address", jsonTag)
		case "min":
			errors[jsonTag] = fmt.Sprintf("%s must be at least %s characters long", jsonTag, fieldError.Param())
		case "eqfield":
			errors[jsonTag] = fmt.Sprintf("%s must match %s", jsonTag, fieldError.Param())
		default:
			errors[jsonTag] = fieldError.Error() // Default error message if not customized
		}

	}

	return errors
}
