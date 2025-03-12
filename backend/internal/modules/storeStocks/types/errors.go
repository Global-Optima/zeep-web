package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var ErrStockAlreadyExists = moduleErrors.NewModuleError(errors.New("stock already exists"))
