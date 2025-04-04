package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrStoreProductNotFound       = moduleErrors.NewModuleError(errors.New("store product not found"))
	ErrStoreProductSizeNotFound   = moduleErrors.NewModuleError(errors.New("store product size not found"))
	ErrInvalidStoreProductID      = moduleErrors.NewModuleError(errors.New("invalid store product ID"))
	ErrStoreProductAlreadyExists  = moduleErrors.NewModuleError(errors.New("store product already exists"))
	ErrInappropriateProductSizeID = moduleErrors.NewModuleError(errors.New("product size ID does not match the product given"))
	ErrValidationFailed           = moduleErrors.NewModuleError(errors.New("validation failed"))
	ErrNoUpdateContext            = moduleErrors.NewModuleError(errors.New("no update context found"))

	ErrStoreStockNotFound            = moduleErrors.NewModuleError(errors.New("store stock not found"))
	ErrInsufficientStock             = moduleErrors.NewModuleError(errors.New("insufficient stock"))
	ErrFailedToFetchStoreProductSize = moduleErrors.NewModuleError(errors.New("failed to fetch store product size"))
	ErrFailedToFetchStoreStock       = moduleErrors.NewModuleError(errors.New("failed to fetch store stock for product size"))
)
