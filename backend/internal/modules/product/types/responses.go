package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500ProductCreate = localization.NewResponseKey(500, data.ProductComponent, data.CreateOperation.ToString())
	Response500ProductUpdate = localization.NewResponseKey(500, data.ProductComponent, data.UpdateOperation.ToString())
	Response500ProductDelete = localization.NewResponseKey(500, data.ProductComponent, data.DeleteOperation.ToString())
	Response500ProductGet    = localization.NewResponseKey(500, data.ProductComponent, data.GetOperation.ToString())
	Response404Product       = localization.NewResponseKey(404, data.ProductComponent)
	Response400Product       = localization.NewResponseKey(400, data.ProductComponent)
	Response201Product       = localization.NewResponseKey(201, data.ProductComponent)
	Response200ProductUpdate = localization.NewResponseKey(200, data.ProductComponent, data.UpdateOperation.ToString())
	Response200ProductDelete = localization.NewResponseKey(200, data.ProductComponent, data.DeleteOperation.ToString())

	Response500ProductSizeCreate    = localization.NewResponseKey(500, data.ProductSizeComponent, data.CreateOperation.ToString())
	Response500ProductSizeUpdate    = localization.NewResponseKey(500, data.ProductSizeComponent, data.UpdateOperation.ToString())
	Response500ProductSizeDelete    = localization.NewResponseKey(500, data.ProductSizeComponent, data.DeleteOperation.ToString())
	Response500ProductSizeGet       = localization.NewResponseKey(500, data.ProductSizeComponent, data.GetOperation.ToString())
	Response404ProductSize          = localization.NewResponseKey(404, data.ProductSizeComponent)
	Response400ProductSize          = localization.NewResponseKey(400, data.ProductSizeComponent)
	Response400ProductSizeDuplicate = localization.NewResponseKey(400, data.ProductSizeComponent, "duplicate")
	Response201ProductSize          = localization.NewResponseKey(201, data.ProductSizeComponent)
	Response200ProductSizeUpdate    = localization.NewResponseKey(200, data.ProductSizeComponent, data.UpdateOperation.ToString())
	Response200ProductSizeDelete    = localization.NewResponseKey(200, data.ProductSizeComponent, data.DeleteOperation.ToString())
)
