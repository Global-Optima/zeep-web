package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateRegionAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.RegionComponent)

	UpdateRegionAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.RegionComponent, &UpdateRegionDTO{})

	DeleteRegionAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.RegionComponent)
)
