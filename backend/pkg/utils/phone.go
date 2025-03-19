package utils

import (
	"github.com/nyaruka/phonenumbers"
)

const (
	DEFAULT_PHONE_NUMBER_REGION               = "KZ"
	INTERNATIONAL_COMPACT_PHONE_NUMBER_FORMAT = phonenumbers.E164
)

func IsValidPhone(rawNumber, region string) bool {
	// if region is unknown set to ""
	phoneNumber, err := phonenumbers.Parse(rawNumber, region)
	if err != nil {
		return false
	}

	return phonenumbers.IsValidNumber(phoneNumber)
}

func FormatPhoneInput(rawNumber string) string {
	phoneNumber, err := phonenumbers.Parse(rawNumber, "")
	if err != nil {
		return ""
	}
	return phonenumbers.Format(phoneNumber, INTERNATIONAL_COMPACT_PHONE_NUMBER_FORMAT)
}
