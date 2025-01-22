package types

import (
	"errors"
)

func ValidateExpirationDays(addDays int) error {
	if addDays <= 0 {
		return errors.New("the number of days to extend must be greater than zero")
	}
	return nil
}
