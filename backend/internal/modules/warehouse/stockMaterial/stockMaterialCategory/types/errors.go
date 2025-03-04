package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrStockMaterialCategoryNotFound         = moduleErrors.NewModuleError(errors.New("stock material category not found"))
	ErrFailedCreateStockMaterialCategory     = moduleErrors.NewModuleError(errors.New("failed to create stock material category"))
	ErrFailedRetrieveStockMaterialCategory   = moduleErrors.NewModuleError(errors.New("failed to fetch stock material category"))
	ErrFailedRetrieveStockMaterialCategories = moduleErrors.NewModuleError(errors.New("failed to fetch stock material categories"))
	ErrFailedUpdateStockMaterialCategory     = moduleErrors.NewModuleError(errors.New("failed to update stock material category"))
	ErrFailedDeleteStockMaterialCategory     = moduleErrors.NewModuleError(errors.New("failed to delete stock material category"))
)
