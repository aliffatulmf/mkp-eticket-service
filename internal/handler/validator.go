package handler

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidateStruct validates a struct and returns a formatted error message
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// HandleValidationError handles validation errors and writes appropriate response
func HandleValidationError(w http.ResponseWriter, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errorMessages := make(map[string]string)
		for _, fieldError := range validationErrors {
			fieldName := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorMessages[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "min":
				errorMessages[fieldName] = fmt.Sprintf("%s must be at least %s characters", fieldName, fieldError.Param())
			case "max":
				errorMessages[fieldName] = fmt.Sprintf("%s must be at most %s characters", fieldName, fieldError.Param())
			default:
				errorMessages[fieldName] = fmt.Sprintf("%s is invalid", fieldName)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf(`{"errors": %v}`, errorMessages)))
		return
	}

	// For other types of errors
	http.Error(w, "Invalid request", http.StatusBadRequest)
}
