package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrValidation        = moduleErrors.NewModuleError(errors.New("input validation failed"))
	ErrIngredientIsInUse = moduleErrors.NewModuleError(errors.New("ingredient is in use"))
)
