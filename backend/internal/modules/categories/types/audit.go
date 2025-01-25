package types

import (
	"github.com/Global-Optima/zeep-web/backend/internal/data"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit/shared"
)

var (
	CreateProductCategoryAuditFactory = shared.NewAuditActionBaseFactory(
		data.CreateOperation, data.ProductCategoryComponent)

	UpdateProductCategoryAuditFactory = shared.NewAuditActionExtendedFactory(
		data.UpdateOperation, data.ProductCategoryComponent, &UpdateProductCategoryDTO{})

	DeleteProductCategoryAuditFactory = shared.NewAuditActionBaseFactory(
		data.DeleteOperation, data.ProductCategoryComponent)
)
