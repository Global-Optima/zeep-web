package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrFailedListStores        = moduleErrors.NewModuleError(errors.New("failed to list stores"))
	ErrFailedCreateWarehouse   = moduleErrors.NewModuleError(errors.New("failed to create warehouse"))
	ErrFailedToFetchWarehouse  = moduleErrors.NewModuleError(errors.New("failed to fetch warehouse"))
	ErrFailedToFetchWarehouses = moduleErrors.NewModuleError(errors.New("failed to fetch warehouses"))
	ErrFailedUpdateWarehouse   = moduleErrors.NewModuleError(errors.New("failed to update warehouse"))
	ErrFailedDeleteWarehouse   = moduleErrors.NewModuleError(errors.New("failed to delete warehouse"))
	ErrWarehouseNotFound       = moduleErrors.NewModuleError(errors.New("warehouse not found"))
)
