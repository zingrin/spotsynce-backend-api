package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using validator tags and returns field-level errors.
func ValidateStruct(s interface{}) map[string]string {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	validationErrors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		field := strings.ToLower(e.Field())
		validationErrors[field] = formatValidationError(e)
	}
	return validationErrors
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", e.Field(), e.Param())
	case "gtfield":
		return fmt.Sprintf("%s must be after %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s failed validation on '%s' tag", e.Field(), e.Tag())
	}
}

// GetValidator returns the shared validator instance.
func GetValidator() *validator.Validate {
	return validate
}
