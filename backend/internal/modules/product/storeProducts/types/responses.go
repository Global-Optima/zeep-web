package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500StoreProduct            = localization.NewResponseKey(500, data.StoreProductComponent)
	Response409StoreProductInUse       = localization.NewResponseKey(409, data.StoreProductComponent, data.DeleteOperation.ToString(), "in_use")
	Response409StoreProductInUseUpdate = localization.NewResponseKey(409, data.StoreProductComponent, data.UpdateOperation.ToString(), "in_use")
	Response404StoreProduct            = localization.NewResponseKey(404, data.StoreProductComponent)
	Response400StoreProduct            = localization.NewResponseKey(400, data.StoreProductComponent)
	Response200StoreProductUpdate      = localization.NewResponseKey(200, data.StoreProductComponent, data.UpdateOperation.ToString())
	Response200StoreProductDelete      = localization.NewResponseKey(200, data.StoreProductComponent, data.DeleteOperation.ToString())
	Response201StoreProduct            = localization.NewResponseKey(201, data.StoreProductComponent)
	Response201StoreProductMultiple    = localization.NewResponseKey(201, data.StoreProductComponent, "multiple")
)
