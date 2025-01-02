package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/handlerErrors"
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
	"net/http"
)

var (
	ErrInvalidProductID = handlerErrors.NewHandlerError(errors.New("invalid product ID"), http.StatusBadRequest)
)

var (
	ErrNothingToUpdate = moduleErrors.NewModuleError(errors.New("nothing to update"))
)
