package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrNotSynchronizedAdditivesNotFound              = moduleErrors.NewModuleError(errors.New("No not synchronized additives"))
	ErrNotSynchronizedProductSizeIngredientsNotFound = moduleErrors.NewModuleError(errors.New("No product size ingredients to synchronize"))
	ErrNotSynchronizedProductSizeAdditivesNotFound   = moduleErrors.NewModuleError(errors.New("No product size additives to synchronize"))
)
