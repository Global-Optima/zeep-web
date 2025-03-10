package utils

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"golang.org/x/text/unicode/norm"
)

// SanitizeString trims, normalizes, and validates user input
func SanitizeString(input string) (string, bool) {
	// Ensure valid UTF-8 encoding
	if !utf8.ValidString(input) {
		return "", false
	}

	// Normalize text to avoid mixed encoding issues (NFC - Canonical Composition)
	normalizedUTF8 := norm.NFC.String(input)

	// Remove leading and trailing spaces
	trimmed := strings.TrimSpace(normalizedUTF8)

	// Replace multiple spaces, tabs, and newlines with a single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized := spaceRegex.ReplaceAllString(trimmed, " ")

	// Remove non-printable/invisible characters (e.g., zero-width space)
	normalized = removeInvisibleCharacters(normalized)

	// Validate final result: Reject empty or meaningless input
	if len(normalized) == 0 {
		return "", false
	}

	return normalized, true
}

// removeInvisibleCharacters removes non-printable, zero-width, and control characters
func removeInvisibleCharacters(input string) string {
	var builder strings.Builder
	for _, char := range input {
		if unicode.IsPrint(char) && !unicode.IsControl(char) {
			builder.WriteRune(char)
		}
	}
	return builder.String()
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
	err := validate.RegisterValidation("customSanitize", ValidateSanitizedString)
	if err != nil {
		fmt.Printf("failed to register custom validator: %v", err)
		return
	}
}

func InitValidators() {
	validate := validator.New()
	RegisterCustomValidators(validate)
}

func RoundToOneDecimal(value float64) float64 {
	return math.Round(value*10) / 10
}
