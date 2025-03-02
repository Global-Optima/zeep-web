package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func SanitizeString(input string) (string, bool) {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) == 0 {
		return "", false
	}
	return trimmed, true
}

func ValidateSanitizedString(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	_, valid := SanitizeString(value)
	return valid
}

func RegisterCustomValidators(validate *validator.Validate) {
	validate.RegisterValidation("customSanitize", ValidateSanitizedString)
}

func InitValidators() {
	validate := validator.New()
	RegisterCustomValidators(validate)
}
