package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrCategoryNotFound = moduleErrors.NewModuleError(errors.New("category not found"))
	ErrCategoryIsInUse  = moduleErrors.NewModuleError(errors.New("category is in use"))
)
