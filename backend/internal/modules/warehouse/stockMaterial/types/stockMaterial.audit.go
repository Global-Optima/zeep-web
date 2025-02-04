package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStockMaterialAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.StockMaterialComponent)

	UpdateStockMaterialAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.StockMaterialComponent, &UpdateStockMaterialDTO{})

	DeleteStockMaterialAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.StockMaterialComponent)
)
