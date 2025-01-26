package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateFranchiseeAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.FranchiseeComponent)

	UpdateFranchiseeAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.FranchiseeComponent, &UpdateFranchiseeDTO{})

	DeleteFranchiseeAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.FranchiseeComponent)
)
