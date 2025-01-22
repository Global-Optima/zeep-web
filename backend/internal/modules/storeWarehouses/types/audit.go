package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStoreWarehouseStockAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.CreateOperation, data.StoreWarehouseStockComponent, &AddStoreStockDTO{})

	UpdateStoreWarehouseStockAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreWarehouseStockComponent, &UpdateStoreStockDTO{})

	DeleteStoreWarehouseStockAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreWarehouseStockComponent, struct{}{})
)
