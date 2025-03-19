package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateUnitAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.UnitComponent)

	UpdateUnitAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.UnitComponent, &UpdateUnitDTO{})

	DeleteUnitAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.UnitComponent)
)
