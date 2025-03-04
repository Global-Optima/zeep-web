package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	// 500 Internal Server Error responses
	Response500StoreCreate = localization.NewResponseKey(500, data.StoreComponent, data.CreateOperation.ToString())
	Response500StoreGet    = localization.NewResponseKey(500, data.StoreComponent, data.GetOperation.ToString())
	Response500StoreUpdate = localization.NewResponseKey(500, data.StoreComponent, data.UpdateOperation.ToString())
	Response500StoreDelete = localization.NewResponseKey(500, data.StoreComponent, data.DeleteOperation.ToString())
	// 404 Not Found response
	Response404Store = localization.NewResponseKey(404, data.StoreComponent)
	// 400 Bad Request responses
	Response400Store = localization.NewResponseKey(400, data.StoreComponent, data.CreateOperation.ToString())
	// 201 Created response
	Response201Store = localization.NewResponseKey(201, data.StoreComponent)
	// 200 Success responses
	Response200StoreUpdate = localization.NewResponseKey(200, data.StoreComponent, data.UpdateOperation.ToString())
	Response200StoreDelete = localization.NewResponseKey(200, data.StoreComponent, data.DeleteOperation.ToString())
)
