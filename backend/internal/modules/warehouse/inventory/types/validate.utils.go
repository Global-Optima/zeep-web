package types

import (
	"errors"
	"time"
)

func ValidateExpirationDate(newExpirationDate, oldExpirationDate time.Time) error {
	if newExpirationDate.Before(oldExpirationDate) {
		return errors.New("new expiration date cannot be earlier than the current expiration date")
	}
	return nil
}
