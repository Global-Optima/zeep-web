package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrValidation                  = moduleErrors.NewModuleError(errors.New("validation error"))
	ErrInvalidWeekdayFormat        = moduleErrors.NewModuleError(errors.New("invalid weekday format"))
	ErrInvalidEmployeeRole         = moduleErrors.NewModuleError(errors.New("invalid employee role"))
	ErrEmployeeAlreadyExists       = moduleErrors.NewModuleError(errors.New("employee already exists"))
	ErrWorkdayAlreadyExists        = moduleErrors.NewModuleError(errors.New("workday already exists"))
	ErrUnsupportedEmployeeType     = moduleErrors.NewModuleError(errors.New("unsupported employee type"))
	ErrNothingToUpdate             = moduleErrors.NewModuleError(errors.New("nothing to update"))
	ErrEmployeeTypeAndRoleMismatch = moduleErrors.NewModuleError(errors.New("employee type and role mismatch"))
)
