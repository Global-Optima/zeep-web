package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrFailedListStores  = moduleErrors.NewModuleError(errors.New("failed to list stores"))
	ErrFailedCreateStore = moduleErrors.NewModuleError(errors.New("failed to create store"))
	ErrFailedFetchStore  = moduleErrors.NewModuleError(errors.New("failed to fetch store"))
	ErrFailedUpdateStore = moduleErrors.NewModuleError(errors.New("failed to update store"))
	ErrFailedDeleteStore = moduleErrors.NewModuleError(errors.New("failed to delete store"))
	ErrStoreNotFound     = moduleErrors.NewModuleError(errors.New("store not found"))
)
