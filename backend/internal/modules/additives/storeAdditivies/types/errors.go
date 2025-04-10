package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrStoreAdditiveAlreadyExists      = moduleErrors.NewModuleError(errors.New("Store additive already exists"))
	ErrStoreAdditiveCategoriesNotFound = moduleErrors.NewModuleError(errors.New("Store additive categories not found"))
	ErrStoreAdditiveNotFound           = moduleErrors.NewModuleError(errors.New("Store additive not found"))
	ErrStoreAdditiveInUse              = moduleErrors.NewModuleError(errors.New("Store additive is in use and cannot be deleted"))

	ErrStoreStockNotFound         = moduleErrors.NewModuleError(errors.New("store stock not found"))
	ErrInsufficientStock          = moduleErrors.NewModuleError(errors.New("insufficient stock"))
	ErrFailedToFetchStoreAdditive = moduleErrors.NewModuleError(errors.New("failed to fetch store additive"))
	ErrFailedToFetchStoreStock    = moduleErrors.NewModuleError(errors.New("failed to fetch store stock for additive"))
)
