package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrSupplierNotFound        = moduleErrors.NewModuleError(errors.New("supplier not found"))
	ErrSupplierExists          = moduleErrors.NewModuleError(errors.New("a supplier with this contact phone already exists"))
	ErrFailedCreateSupplier    = moduleErrors.NewModuleError(errors.New("failed to create supplier"))
	ErrFailedRetrieveSupplier  = moduleErrors.NewModuleError(errors.New("failed to fetch supplier"))
	ErrFailedListSupplier      = moduleErrors.NewModuleError(errors.New("failed to list suppliers"))
	ErrFailedUpdateSupplier    = moduleErrors.NewModuleError(errors.New("failed to update supplier"))
	ErrFailedDeleteSupplier    = moduleErrors.NewModuleError(errors.New("failed to delete supplier"))
	ErrFailedToFetchMaterials  = moduleErrors.NewModuleError(errors.New("failed to fetch supplier materials"))
	ErrFailedToUpsertMaterials = moduleErrors.NewModuleError(errors.New("failed to upsert supplier materials"))
)
