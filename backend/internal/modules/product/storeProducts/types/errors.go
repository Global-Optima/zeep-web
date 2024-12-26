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
	ErrStoreProductNotFound       = moduleErrors.NewModuleError(errors.New("store product not found"))
	ErrStoreProductSizeNotFound   = moduleErrors.NewModuleError(errors.New("store product size not found"))
	ErrInvalidStoreProductID      = moduleErrors.NewModuleError(errors.New("invalid store product ID"))
	ErrStoreProductAlreadyExists  = moduleErrors.NewModuleError(errors.New("store product already exists"))
	ErrInappropriateProductSizeID = moduleErrors.NewModuleError(errors.New("product size ID does not match the product given"))
	ErrValidationFailed           = moduleErrors.NewModuleError(errors.New("validation failed"))
)
