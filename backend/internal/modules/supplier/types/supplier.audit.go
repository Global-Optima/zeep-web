package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateSupplierAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.SupplierComponent)

	UpdateSupplierAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.SupplierComponent, &UpdateSupplierDTO{})

	DeleteSupplierAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.SupplierComponent)
)
