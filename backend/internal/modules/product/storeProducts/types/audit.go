package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateStoreProductAuditFactory = shared.NewAuditActionExtendedFactory(
		data.CreateOperation, data.StoreProductComponent, &AuditStoreProductDTO{})

	UpdateStoreProductAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.StoreProductComponent, &UpdateStoreProductDTO{})

	DeleteStoreProductAuditFactory = shared.NewAuditActionExtendedFactory(
		data.DeleteOperation, data.StoreProductComponent, &AuditStoreProductDTO{})
)
