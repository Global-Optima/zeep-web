package moduleErrors

import "github.com/pkg/errors"

var (
	ErrNotFound            = NewModuleError(errors.New("resource not found"))
	ErrInvalidPhoneNumber  = NewModuleError(errors.New("invalid phone number"))
	ErrInvalidEmailAddress = NewModuleError(errors.New("invalid email address"))
	ErrInvalidHoursFormat  = NewModuleError(errors.New("invalid hours format"))
	ErrInvalidPassword     = NewModuleError(errors.New("invalid password"))
	ErrInvalidWeekday      = NewModuleError(errors.New("invalid weekday"))

	ValidationErrorsMap = map[*ModuleError]string{
		ErrInvalidPhoneNumber:  "phone",
		ErrInvalidEmailAddress: "email",
		ErrInvalidHoursFormat:  "hours",
		ErrInvalidPassword:     "password",
		ErrInvalidWeekday:      "weekday",
	}
)
