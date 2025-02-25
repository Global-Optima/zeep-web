package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrInactiveEmployee   = moduleErrors.NewModuleError(errors.New("inactive employee"))
	ErrBannedCustomer     = moduleErrors.NewModuleError(errors.New("banned customer"))
	ErrInvalidCredentials = moduleErrors.NewModuleError(errors.New("invalid credentials"))
)
