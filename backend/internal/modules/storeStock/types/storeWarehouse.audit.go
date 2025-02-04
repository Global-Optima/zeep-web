package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStoreStockAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.CreateOperation, data.StoreStockComponent, &AddStoreStockDTO{})

	UpdateStoreStockAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreStockComponent, &UpdateStoreStockDTO{})

	DeleteStoreStockAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreStockComponent, struct{}{})
)
