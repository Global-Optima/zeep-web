package types

import (
	"net/http"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

// Handler errors
var (
	ErrInvalidUnitID       = handlerErrors.NewHandlerError(errors.New("invalid unit ID"), http.StatusBadRequest)
	ErrInvalidBody         = handlerErrors.NewHandlerError(errors.New("invalid body"), http.StatusBadRequest)
	ErrFailedToParseFilter = handlerErrors.NewHandlerError(errors.New("failed to parse filter parameters"), http.StatusBadRequest)
	ErrFailedToFetchUnits  = handlerErrors.NewHandlerError(errors.New("failed to fetch units"), http.StatusInternalServerError)
	ErrFailedToCreateUnit  = handlerErrors.NewHandlerError(errors.New("failed to create unit"), http.StatusInternalServerError)
	ErrFailedToUpdateUnit  = handlerErrors.NewHandlerError(errors.New("failed to update unit"), http.StatusInternalServerError)
	ErrFailedToDeleteUnit  = handlerErrors.NewHandlerError(errors.New("failed to delete unit"), http.StatusInternalServerError)
)

// Module errors
var (
	ErrUnitNotFound        = moduleErrors.NewModuleError(errors.New("unit not found"))
	ErrNothingToUpdate     = moduleErrors.NewModuleError(errors.New("nothing to update"))
	ErrFailedToApplyFilter = moduleErrors.NewModuleError(errors.New("failed to apply filter"))
)
