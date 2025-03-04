package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/localization"
)

var (
	// 400 Bad Request responses
	Response400StockMaterialCategory = localization.NewResponseKey(400, data.StockMaterialCategoryComponent, data.CreateOperation.ToString())

	// 200 Success responses
	Response200StockMaterialCategoryUpdate = localization.NewResponseKey(200, data.StockMaterialCategoryComponent, data.UpdateOperation.ToString())
	Response200StockMaterialCategoryDelete = localization.NewResponseKey(200, data.StockMaterialCategoryComponent, data.DeleteOperation.ToString())

	// 201 Created response
	Response201StockMaterialCategory = localization.NewResponseKey(201, data.StockMaterialCategoryComponent)

	// 500 Internal Server Error responses
	Response500StockMaterialCategoryCreate = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.CreateOperation.ToString())
	Response500StockMaterialCategoryGet    = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.GetOperation.ToString())
	Response500StockMaterialCategoryUpdate = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.UpdateOperation.ToString())
	Response500StockMaterialCategoryDelete = localization.NewResponseKey(500, data.StockMaterialCategoryComponent, data.DeleteOperation.ToString())

	// 404 Not Found response
	Response404StockMaterialCategory = localization.NewResponseKey(404, data.StockMaterialCategoryComponent)
)
