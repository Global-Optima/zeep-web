package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStoreAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.StoreComponent)

	UpdateStoreAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.StoreComponent, &UpdateStoreDTO{})

	DeleteStoreAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.StoreComponent)
)
