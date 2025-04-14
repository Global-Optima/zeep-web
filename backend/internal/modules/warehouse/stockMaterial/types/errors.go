package types

import (
	"errors"

	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
)

var (
	ErrStockMaterialNotFound        = moduleErrors.NewModuleError(errors.New("stock material not found"))
	ErrStockMaterialBarcodeNotFound = moduleErrors.NewModuleError(errors.New("stockMaterial not found with the provided barcode"))
	ErrStockMaterialInUse           = moduleErrors.NewModuleError(errors.New("stock material is in use"))
)
