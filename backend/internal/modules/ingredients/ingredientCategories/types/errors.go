package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrIngredientCategoryNotFound      = moduleErrors.NewModuleError(errors.New("ingredient category not found"))
	ErrIngredientCategoryIsInUse       = moduleErrors.NewModuleError(errors.New("ingredient category is in use"))
	ErrFailedToFetchIngredientCategory = moduleErrors.NewModuleError(errors.New("failed to fetch ingredient category"))
)
