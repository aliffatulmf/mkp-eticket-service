package validator

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s any) error {
	return validate.Struct(s)
}

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
		fmt.Fprintf(w, `{"errors": %v}`, errorMessages)
		return
	}

	http.Error(w, "Invalid request", http.StatusBadRequest)
}
