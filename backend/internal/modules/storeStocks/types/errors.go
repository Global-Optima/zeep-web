package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrStockAlreadyExists = moduleErrors.NewModuleError(errors.New("stock already exists"))
	ErrStockNotFound      = moduleErrors.NewModuleError(errors.New("stock not found"))
	ErrStockIsInUse       = moduleErrors.NewModuleError(errors.New("stock is in use"))
	ErrInsufficientStock  = moduleErrors.NewModuleError(errors.New("insufficient stock"))
)
