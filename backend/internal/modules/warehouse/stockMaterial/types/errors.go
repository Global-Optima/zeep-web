package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var ErrStockMaterialNotFound = moduleErrors.NewModuleError(errors.New("stock material not found"))
