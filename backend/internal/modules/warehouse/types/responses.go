package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500WarehouseCreate = localization.NewResponseKey(500, data.WarehouseComponent, data.CreateOperation.ToString())
	Response500WarehouseUpdate = localization.NewResponseKey(500, data.WarehouseComponent, data.UpdateOperation.ToString())
	Response500WarehouseGet    = localization.NewResponseKey(500, data.WarehouseComponent, data.GetOperation.ToString())
	Response500WarehouseDelete = localization.NewResponseKey(500, data.WarehouseComponent, data.DeleteOperation.ToString())
	Response500WarehouseAssign = localization.NewResponseKey(500, data.WarehouseComponent, "ASSIGN")
	Response404Warehouse       = localization.NewResponseKey(404, data.WarehouseComponent)
	Response400Warehouse       = localization.NewResponseKey(400, data.WarehouseComponent)
	Response200WarehouseUpdate = localization.NewResponseKey(200, data.WarehouseComponent, data.UpdateOperation.ToString())
	Response200WarehouseDelete = localization.NewResponseKey(200, data.WarehouseComponent, data.DeleteOperation.ToString())
	Response201Warehouse       = localization.NewResponseKey(201, data.WarehouseComponent)
	Response200WarehouseAssign = localization.NewResponseKey(200, data.WarehouseComponent, "ASSIGN")
)
