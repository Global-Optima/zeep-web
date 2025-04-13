package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	ErrProductAlreadyExists = moduleErrors.NewModuleError(errors.New("product already exists"))
	ErrProductNotFound      = moduleErrors.NewModuleError(errors.New("product not found"))
	ErrProductIsInUse       = moduleErrors.NewModuleError(errors.New("product in use"))

	ErrProductSizeIngredientsNotFound      = moduleErrors.NewModuleError(errors.New("product size ingredients not found"))
	ErrProductSizeDefaultAdditivesNotFound = moduleErrors.NewModuleError(errors.New("product size default additives not found"))
	ErrProductSizeNotFound                 = moduleErrors.NewModuleError(errors.New("product size not found"))
	ErrProductSizeUniqueName               = moduleErrors.NewModuleError(errors.New("product size unique name"))
	ErrProductSizeIsInUse                  = moduleErrors.NewModuleError(errors.New("product size in use"))
)
