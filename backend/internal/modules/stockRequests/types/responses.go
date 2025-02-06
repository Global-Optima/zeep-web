package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500StockRequest       = localization.NewResponseKey(500, data.StockRequestComponent)
	Response403StockRequest       = localization.NewResponseKey(403, data.StockRequestComponent)
	Response400StockRequest       = localization.NewResponseKey(400, data.StockRequestComponent)
	Response200StockRequestUpdate = localization.NewResponseKey(200, data.StockRequestComponent, data.UpdateOperation.ToString())
	Response200StockRequestDelete = localization.NewResponseKey(200, data.StockRequestComponent, data.DeleteOperation.ToString())
	Response201StockRequest       = localization.NewResponseKey(201, data.StockRequestComponent)
)
