package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response500ProductCategory       = localization.NewResponseKey(500, data.ProductCategoryComponent)
	Response400ProductCategory       = localization.NewResponseKey(400, data.ProductCategoryComponent)
	Response404ProductCategory       = localization.NewResponseKey(404, data.ProductCategoryComponent)
	Response200ProductCategoryUpdate = localization.NewResponseKey(200, data.ProductCategoryComponent, data.UpdateOperation.ToString())
	Response200ProductCategoryDelete = localization.NewResponseKey(200, data.ProductCategoryComponent, data.DeleteOperation.ToString())
	Response201ProductCategory       = localization.NewResponseKey(201, data.ProductCategoryComponent)
)
