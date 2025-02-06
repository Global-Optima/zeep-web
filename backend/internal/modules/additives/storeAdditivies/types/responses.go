package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500StoreAdditive       = localization.NewResponseKey(500, data.StoreAdditiveComponent)
	Response400StoreAdditive       = localization.NewResponseKey(400, data.StoreAdditiveComponent)
	Response200StoreAdditiveUpdate = localization.NewResponseKey(200, data.StoreAdditiveComponent, data.UpdateOperation.ToString())
	Response200StoreAdditiveDelete = localization.NewResponseKey(200, data.StoreAdditiveComponent, data.DeleteOperation.ToString())
	Response201StoreAdditive       = localization.NewResponseKey(201, data.StoreAdditiveComponent)
)
