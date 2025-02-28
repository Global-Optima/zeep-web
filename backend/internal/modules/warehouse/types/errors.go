package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var ErrWarehouseNotFound = moduleErrors.NewModuleError(errors.New("warehouse not found"))
