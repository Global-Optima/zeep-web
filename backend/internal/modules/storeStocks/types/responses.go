package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500StoreStock         = localization.NewResponseKey(500, data.StoreStockComponent)
	Response400StoreStock         = localization.NewResponseKey(400, data.StoreStockComponent)
	Response200StoreStockUpdate   = localization.NewResponseKey(200, data.StoreStockComponent, data.UpdateOperation.ToString())
	Response200StoreStockDelete   = localization.NewResponseKey(200, data.StoreStockComponent, data.DeleteOperation.ToString())
	Response201StoreStock         = localization.NewResponseKey(201, data.StoreStockComponent)
	Response201StoreStockMultiple = localization.NewResponseKey(201, data.StoreStockComponent, "multiple")
)
