package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response400StockMaterial                = localization.NewResponseKey(400, data.StockMaterialComponent)
	Response500StockMaterialGet             = localization.NewResponseKey(500, data.StockMaterialComponent, data.GetOperation.ToString())
	Response500StockMaterialCreate          = localization.NewResponseKey(500, data.StockMaterialComponent, data.CreateOperation.ToString())
	Response500StockMaterialUpdate          = localization.NewResponseKey(500, data.StockMaterialComponent, data.UpdateOperation.ToString())
	Response500StockMaterialDelete          = localization.NewResponseKey(500, data.StockMaterialComponent, data.DeleteOperation.ToString())
	Response404StockMaterial                = localization.NewResponseKey(404, data.StockMaterialComponent)
	Response500StockMaterialDeactivate      = localization.NewResponseKey(500, data.StockMaterialComponent, "deactivate")
	Response500StockMaterialBarcode         = localization.NewResponseKey(500, data.StockMaterialComponent, "barcode")
	Response500StockMaterialBarcodeGet      = localization.NewResponseKey(500, data.StockMaterialComponent, "barcode-get")
	Response400StockMaterialBarcodeRequired = localization.NewResponseKey(400, data.StockMaterialComponent, "barcode-required")
	Response404StockMaterialBarcode         = localization.NewResponseKey(404, data.StockMaterialComponent, "barcode")
)
