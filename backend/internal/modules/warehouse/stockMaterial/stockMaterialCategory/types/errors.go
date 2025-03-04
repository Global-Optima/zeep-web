package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrStockMaterialCategoryNotFound = moduleErrors.NewModuleError(errors.New("stock material category not found"))
	ErrStockMaterialCategoryIsInUse  = moduleErrors.NewModuleError(errors.New("stock material category is in use"))
)
