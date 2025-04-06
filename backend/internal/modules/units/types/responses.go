package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500UnitCreate = localization.NewResponseKey(500, data.UnitComponent, data.CreateOperation.ToString())
	Response500UnitGet    = localization.NewResponseKey(500, data.UnitComponent, data.GetOperation.ToString())
	Response500UnitUpdate = localization.NewResponseKey(500, data.UnitComponent, data.UpdateOperation.ToString())
	Response500UnitDelete = localization.NewResponseKey(500, data.UnitComponent, data.DeleteOperation.ToString())

	Response409UnitDeleteInUse = localization.NewResponseKey(409, data.UnitComponent, data.DeleteOperation.ToString(), "in-use")
	Response404Unit            = localization.NewResponseKey(404, data.UnitComponent)
	Response400Unit            = localization.NewResponseKey(400, data.UnitComponent)

	Response201Unit       = localization.NewResponseKey(201, data.UnitComponent)
	Response200UnitUpdate = localization.NewResponseKey(200, data.UnitComponent, data.UpdateOperation.ToString())
	Response200UnitDelete = localization.NewResponseKey(200, data.UnitComponent, data.DeleteOperation.ToString())
)
