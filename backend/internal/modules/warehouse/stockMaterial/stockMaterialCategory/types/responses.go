package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	Response400StockMaterialCategory            = localization.NewResponseKey(400, data.StockMaterialCategoryComponent)
	Response500StockMaterialCategoryCreate      = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.CreateOperation.ToString())
	Response500StockMaterialCategoryGet         = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.GetOperation.ToString())
	Response500StockMaterialCategoryUpdate      = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.UpdateOperation.ToString())
	Response500StockMaterialCategoryDelete      = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.DeleteOperation.ToString())
	Response404StockMaterialCategory            = localization.NewResponseKey(404, data.StockMaterialCategoryComponent)
	Response409StockMaterialCategoryDeleteInUse = localization.NewResponseKey(409, data.StockMaterialCategoryComponent, data.DeleteOperation.ToString(), "in-use")
	Response200StockMaterialCategoryUpdate      = localization.NewResponseKey(200, data.StockMaterialCategoryComponent, data.UpdateOperation.ToString())
	Response200StockMaterialCategoryDelete      = localization.NewResponseKey(200, data.StockMaterialCategoryComponent, data.DeleteOperation.ToString())
	Response201StockMaterialCategoryCreate      = localization.NewResponseKey(201, data.StockMaterialCategoryComponent, data.CreateOperation.ToString())
)
