package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/warehouse/stockMaterial/stockMaterialCategory"
)

type StockMaterialCategoriesModule struct {
	*common.BaseModule
	Repo    stockMaterialCategory.StockMaterialCategoryRepository
	Service stockMaterialCategory.StockMaterialCategoryService
	Handler *stockMaterialCategory.StockMaterialCategoryHandler
}

func NewStockMaterialCategoriesModule(base *common.BaseModule, auditService audit.AuditService) *StockMaterialCategoriesModule {
	repo := stockMaterialCategory.NewStockMaterialCategoryRepository(base.DB)
	service := stockMaterialCategory.NewStockMaterialCategoryService(repo, base.Logger)
	handler := stockMaterialCategory.NewStockMaterialCategoryHandler(service, auditService)

	base.Router.RegisterStockMaterialCategoryRoutes(handler)

	return &StockMaterialCategoriesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
