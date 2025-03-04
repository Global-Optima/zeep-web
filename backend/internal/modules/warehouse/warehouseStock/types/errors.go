package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/errors/moduleErrors"
	"github.com/pkg/errors"
)

var (
	// Inventory processing errors
	ErrReceiveInventory = moduleErrors.NewModuleError(errors.New("failed to receive inventory"))

	// Delivery errors
	ErrFetchDeliveries          = moduleErrors.NewModuleError(errors.New("failed to fetch deliveries"))
	ErrFetchDelivery            = moduleErrors.NewModuleError(errors.New("failed to fetch delivery"))
	ErrFailedToRecordDeliveries = moduleErrors.NewModuleError(errors.New("failed to record deliveries"))
	ErrDeliveryNotFound         = moduleErrors.NewModuleError(errors.New("delivery not found"))

	// Stock material update errors
	ErrAddWarehouseStockMaterial = moduleErrors.NewModuleError(errors.New("failed to add warehouse stock material"))
	ErrDeductFromStock           = moduleErrors.NewModuleError(errors.New("failed to deduct from stock"))
	ErrEmptyStocks               = moduleErrors.NewModuleError(errors.New("stocks cannot be empty"))
	ErrUpdateStock               = moduleErrors.NewModuleError(errors.New("failed to update warehouse stock"))
	ErrUpdateStockQuantity       = moduleErrors.NewModuleError(errors.New("failed to update warehouse stock quantity"))
	ErrUpdateExpiration          = moduleErrors.NewModuleError(errors.New("failed to update expiration date"))
	ErrNothingToUpdate           = moduleErrors.NewModuleError(errors.New("nothing to update"))

	// Stock retrieval errors
	ErrFetchStock                = moduleErrors.NewModuleError(errors.New("failed to fetch warehouse stocks"))
	ErrFetchStockMaterials       = moduleErrors.NewModuleError(errors.New("failed to fetch stock materials"))
	ErrFetchStockMaterialDetails = moduleErrors.NewModuleError(errors.New("failed to fetch stock material details"))

	// Not found errors
	ErrStockMaterialNotFound = moduleErrors.NewModuleError(errors.New("stock material not found"))
)
