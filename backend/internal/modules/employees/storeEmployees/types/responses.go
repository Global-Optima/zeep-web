package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500StoreEmployeeGet    = localization.NewResponseKey(500, data.StoreEmployeeComponent, data.GetOperation.ToString())
	Response500StoreEmployeeCreate = localization.NewResponseKey(500, data.StoreEmployeeComponent, data.CreateOperation.ToString())
	Response500StoreEmployeeUpdate = localization.NewResponseKey(500, data.StoreEmployeeComponent, data.UpdateOperation.ToString())
	Response500StoreEmployeeDelete = localization.NewResponseKey(500, data.StoreEmployeeComponent, data.DeleteOperation.ToString())

	Response400StoreEmployee = localization.NewResponseKey(400, data.StoreEmployeeComponent)

	Response200StoreEmployeeUpdate = localization.NewResponseKey(200, data.StoreEmployeeComponent, data.UpdateOperation.ToString())
	Response200StoreEmployeeDelete = localization.NewResponseKey(200, data.StoreEmployeeComponent, data.DeleteOperation.ToString())
	Response201StoreEmployee       = localization.NewResponseKey(201, data.StoreEmployeeComponent)
)
