package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	// 400 Bad Request responses
	Response400Region = localization.NewResponseKey(400, data.RegionComponent, data.CreateOperation.ToString())

	// 200 Success responses
	Response200RegionUpdate = localization.NewResponseKey(200, data.RegionComponent, data.UpdateOperation.ToString())
	Response200RegionDelete = localization.NewResponseKey(200, data.RegionComponent, data.DeleteOperation.ToString())

	// 201 Created response
	Response201Region = localization.NewResponseKey(201, data.RegionComponent)

	// 500 Internal Server Error responses
	Response500RegionCreate = localization.NewResponseKey(500, data.RegionComponent, data.CreateOperation.ToString())
	Response500RegionUpdate = localization.NewResponseKey(500, data.RegionComponent, data.UpdateOperation.ToString())
	Response500RegionGet    = localization.NewResponseKey(500, data.RegionComponent, data.GetOperation.ToString())
	Response500RegionDelete = localization.NewResponseKey(500, data.RegionComponent, data.DeleteOperation.ToString())

	// 404 Not Found response
	Response404Region = localization.NewResponseKey(404, data.RegionComponent)
)
