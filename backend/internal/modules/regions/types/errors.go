package types

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrFailedToCheckRegionWarehouse = handlerErrors.NewHandlerError(errors.New("failed to check region warehouse"), http.StatusInternalServerError)
	ErrRegionWarehouseMismatch      = handlerErrors.NewHandlerError(errors.New("this warehouse is not assigned to the given region"), http.StatusForbidden)

	ErrFailedCreateRegion = moduleErrors.NewModuleError(errors.New("failed to create region"))
	ErrFailedUpdateRegion = moduleErrors.NewModuleError(errors.New("failed to update region"))
	ErrFailedDeleteRegion = moduleErrors.NewModuleError(errors.New("failed to delete region"))
	ErrFailedFetchRegion  = moduleErrors.NewModuleError(errors.New("failed to fetch region"))
	ErrFailedFetchRegions = moduleErrors.NewModuleError(errors.New("failed to fetch regions"))
	ErrRegionNotFound     = moduleErrors.NewModuleError(errors.New("region not found"))
	ErrFailedCheckRegion  = moduleErrors.NewModuleError(errors.New("failed to check region warehouse association"))
)
