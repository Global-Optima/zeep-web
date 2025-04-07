package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateProvisionAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.ProvisionComponent)

	UpdateProvisionAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.ProvisionComponent, &UpdateProvisionDTO{})

	DeleteProvisionAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.ProvisionComponent)
)
