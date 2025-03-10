package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrInsufficientStock = moduleErrors.NewModuleError(errors.New("insufficient stock to fulfill the order"))
	ErrMultipleSelect    = moduleErrors.NewModuleError(errors.New("multiple select on this additive category is not allowed"))
)
