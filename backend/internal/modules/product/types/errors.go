package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrInvalidStoreID     = handlerErrors.NewHandlerError(errors.New("invalid store ID"), http.StatusBadRequest)
	ErrUnauthorizedAccess = handlerErrors.NewHandlerError(errors.New("unauthorized access to store"), http.StatusUnauthorized)
)

var (
	ErrProductNotFound = moduleErrors.NewModuleError(errors.New("product not found"))
)
