package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateProductAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.ProductComponent)

	CreateProductSizeAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.ProductSizeComponent)

	UpdateProductAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.ProductComponent, &UpdateProductDTO{})

	UpdateProductSizeAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.ProductSizeComponent, &UpdateProductSizeDTO{})

	DeleteProductAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.ProductComponent)

	DeleteProductSizeAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.ProductSizeComponent)
)
