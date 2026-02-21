package middleware

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// InitValidator initializes the validator
func InitValidator() {
	validate = validator.New()
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	return validate
}

// ValidateStruct validates a struct against the defined tags
func ValidateStruct(data interface{}) error {
	if validate == nil {
		InitValidator()
	}
	return validate.Struct(data)
}

// FormatValidationErrors formats validator errors into a map
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, err := range validationErrors {
			errors[err.Field()] = err.Tag()
		}
	}
	return errors
}
