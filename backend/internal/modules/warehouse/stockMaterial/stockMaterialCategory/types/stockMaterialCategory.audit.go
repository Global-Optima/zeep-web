package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStockMaterialAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.StockMaterialCategoryComponent)

	UpdateStockMaterialAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.StockMaterialCategoryComponent, &UpdateStockMaterialCategoryDTO{})

	DeleteStockMaterialAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.StockMaterialCategoryComponent)
)
