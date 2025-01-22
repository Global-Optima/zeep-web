package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
)

type CategoriesModule struct {
	*common.BaseModule
	Repo    categories.CategoryRepository
	Service categories.CategoryService
	Handler *categories.CategoryHandler
}

func NewCategoriesModule(base *common.BaseModule, auditService audit.AuditService) *CategoriesModule {
	repo := categories.NewCategoryRepository(base.DB)
	service := categories.NewCategoryService(repo)
	handler := categories.NewCategoryHandler(service, auditService)

	base.Router.RegisterProductCategoriesRoutes(handler)

	return &CategoriesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
