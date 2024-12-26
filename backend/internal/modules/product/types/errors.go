package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrProductNotFound     = moduleErrors.NewModuleError(errors.New("product not found"))
	ErrProductSizeNotFound = moduleErrors.NewModuleError(errors.New("product size not found"))
)
