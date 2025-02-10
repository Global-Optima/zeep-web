package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500RegionEmployeeGet    = localization.NewResponseKey(500, data.RegionEmployeeComponent, data.GetOperation.ToString())
	Response500RegionEmployeeCreate = localization.NewResponseKey(500, data.RegionEmployeeComponent, data.CreateOperation.ToString())
	Response500RegionEmployeeUpdate = localization.NewResponseKey(500, data.RegionEmployeeComponent, data.UpdateOperation.ToString())
	Response500RegionEmployeeDelete = localization.NewResponseKey(500, data.RegionEmployeeComponent, data.DeleteOperation.ToString())

	Response400RegionEmployee = localization.NewResponseKey(400, data.RegionEmployeeComponent)

	Response200RegionEmployeeUpdate = localization.NewResponseKey(200, data.RegionEmployeeComponent, data.UpdateOperation.ToString())
	Response200RegionEmployeeDelete = localization.NewResponseKey(200, data.RegionEmployeeComponent, data.DeleteOperation.ToString())
	Response201RegionEmployee       = localization.NewResponseKey(201, data.RegionEmployeeComponent)
)
