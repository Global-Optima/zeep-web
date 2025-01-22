package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStoreProductAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.CreateOperation, data.StoreProductComponent, &CreateStoreProductDTO{})

	UpdateStoreProductAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.UpdateOperation, data.StoreProductComponent, &UpdateStoreProductDTO{})

	DeleteStoreProductAuditFactory = shared.NewAuditStoreActionExtendedFactory(
		data.DeleteOperation, data.StoreProductComponent, &struct{}{})
)
