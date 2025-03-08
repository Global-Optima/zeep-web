package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500WarehouseEmployeeGet    = localization.NewResponseKey(500, data.WarehouseEmployeeComponent, data.GetOperation.ToString())
	Response500WarehouseEmployeeCreate = localization.NewResponseKey(500, data.WarehouseEmployeeComponent, data.CreateOperation.ToString())
	Response500WarehouseEmployeeUpdate = localization.NewResponseKey(500, data.WarehouseEmployeeComponent, data.UpdateOperation.ToString())
	Response500WarehouseEmployeeDelete = localization.NewResponseKey(500, data.WarehouseEmployeeComponent, data.DeleteOperation.ToString())

	Response409WarehouseEmployee = localization.NewResponseKey(409, data.WarehouseEmployeeComponent)
	Response400WarehouseEmployee = localization.NewResponseKey(400, data.WarehouseEmployeeComponent)

	Response200WarehouseEmployeeUpdate = localization.NewResponseKey(200, data.WarehouseEmployeeComponent, data.UpdateOperation.ToString())
	Response200WarehouseEmployeeDelete = localization.NewResponseKey(200, data.WarehouseEmployeeComponent, data.DeleteOperation.ToString())
	Response201WarehouseEmployee       = localization.NewResponseKey(201, data.WarehouseEmployeeComponent)
)
