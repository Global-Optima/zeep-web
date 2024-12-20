package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ValidateTime(s string) error {
	re := regexp.MustCompile(`^([0-1][0-9]|2[0-3]):([0-5][0-9])(:([0-5][0-9]))?$`)
	if !re.MatchString(s) {
		return fmt.Errorf("invalid time format")
	}

	parts := strings.Split(s, ":")
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}

	if hours < 0 || hours >= 24 {
		return fmt.Errorf("invalid hours: must be less than 24")
	}

	if minutes < 0 || minutes >= 60 {
		return fmt.Errorf("invalid minutes: must be less than 60")
	}

	if len(parts) == 3 {
		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return err
		}
		if seconds < 0 || seconds >= 60 {
			return fmt.Errorf("invalid seconds: must be less than 60")
		}
	}

	return nil
}
