package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStoreAdditiveAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.CreateOperation, data.StoreAdditiveComponent, &CreateStoreAdditiveDTO{})

	UpdateStoreAdditiveAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreAdditiveComponent, &UpdateStoreAdditiveDTO{})

	DeleteStoreAdditiveAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreAdditiveComponent, struct{}{})
)
