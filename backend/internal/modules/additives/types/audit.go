package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateAdditiveAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.AdditiveComponent)

	CreateAdditiveCategoryAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.AdditiveCategoryComponent)

	UpdateAdditiveAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.AdditiveComponent, &UpdateAdditiveDTO{})

	UpdateAdditiveCategoryAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.AdditiveCategoryComponent, &UpdateAdditiveCategoryDTO{})

	DeleteAdditiveAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.AdditiveComponent)

	DeleteAdditiveCategoryAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.AdditiveCategoryComponent)
)
