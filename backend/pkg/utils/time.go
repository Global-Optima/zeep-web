package utils

import (
	"fmt"
	"regexp"
)

const TIME_REGEXP = `^([0-1][0-9]|2[0-3]):([0-5][0-9])(:([0-5][0-9]))?$`

func ValidateTime(s string) error {
	re := regexp.MustCompile(TIME_REGEXP)
	if !re.MatchString(s) {
		return fmt.Errorf("invalid time format")
	}

	return nil
}
