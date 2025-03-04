package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	// 400 Bad Request responses
	Response400WarehouseStockReceive = localization.NewResponseKey(400, data.WarehouseStockComponent, "RECEIVE")
	Response400WarehouseStockQuery   = localization.NewResponseKey(400, data.WarehouseStockComponent, "QUERY")

	// 200 Success responses
	Response200WarehouseStockReceive     = localization.NewResponseKey(200, data.WarehouseStockComponent, "RECEIVE")
	Response200WarehouseStockAddMaterial = localization.NewResponseKey(200, data.WarehouseStockComponent, "ADD_MATERIAL")
	Response200WarehouseStockDeductStock = localization.NewResponseKey(200, data.WarehouseStockComponent, "DEDUCT_STOCK")
	Response200WarehouseStockUpdate      = localization.NewResponseKey(200, data.WarehouseStockComponent, data.UpdateOperation.ToString())

	// 201 Created response
	Response201WarehouseStock = localization.NewResponseKey(201, data.WarehouseStockComponent)

	// 500 Internal Server Error responses
	Response500WarehouseStockReceive         = localization.NewResponseKey(500, data.WarehouseStockComponent, "RECEIVE")
	Response500WarehouseStockFetchDeliveries = localization.NewResponseKey(500, data.WarehouseStockComponent, "FETCH_DELIVERIES")
	Response500WarehouseStockFetchDelivery   = localization.NewResponseKey(500, data.WarehouseStockComponent, "FETCH_DELIVERY")
	Response500WarehouseStockAddMaterial     = localization.NewResponseKey(500, data.WarehouseStockComponent, "ADD_MATERIAL")
	Response500WarehouseStockDeductStock     = localization.NewResponseKey(500, data.WarehouseStockComponent, "DEDUCT_STOCK")
	Response500WarehouseStockFetchStock      = localization.NewResponseKey(500, data.WarehouseStockComponent, "FETCH_STOCK")
	Response500WarehouseStockFetchDetails    = localization.NewResponseKey(500, data.WarehouseStockComponent, "FETCH_DETAILS")
	Response500WarehouseStockUpdate          = localization.NewResponseKey(500, data.WarehouseStockComponent, data.UpdateOperation.ToString())
	Response500WarehouseStockAddStocks       = localization.NewResponseKey(500, data.WarehouseStockComponent, "ADD_STOCKS")
)
