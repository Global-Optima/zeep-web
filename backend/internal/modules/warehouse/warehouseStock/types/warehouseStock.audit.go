package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

type UpdateAuditDTOs struct {
	*UpdateWarehouseStockDTO
	*AdjustWarehouseStock
}

var (
	CreateStockMaterialAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.CreateOperation, data.WarehouseStockComponent, &AddWarehouseStockMaterial{})

	UpdateStockMaterialAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.UpdateOperation, data.WarehouseStockComponent, &UpdateAuditDTOs{})

	DeleteStockMaterialAuditFactory = shared.NewAuditWarehouseActionExtendedFactory(
		data.DeleteOperation, data.WarehouseStockComponent, struct{}{})
)
