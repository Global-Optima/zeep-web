package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500Product       = localization.NewResponseKey(500, data.ProductComponent)
	Response400Product       = localization.NewResponseKey(400, data.ProductComponent)
	Response200ProductUpdate = localization.NewResponseKey(200, data.ProductComponent, data.UpdateOperation.ToString())
	Response200ProductDelete = localization.NewResponseKey(200, data.ProductComponent, data.DeleteOperation.ToString())
	Response201Product       = localization.NewResponseKey(201, data.ProductComponent)
)
