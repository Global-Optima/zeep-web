package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateWarehouseAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.WarehouseComponent)

	UpdateWarehouseAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.WarehouseComponent, &UpdateWarehouseDTO{})

	DeleteWarehouseAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.WarehouseComponent)
)
