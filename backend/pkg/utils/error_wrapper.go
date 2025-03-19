package utils

import "fmt"

func WrapError(message string, err error) error {
	return fmt.Errorf("%s: %w", message, err)
}

func StringOrEmpty(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}
