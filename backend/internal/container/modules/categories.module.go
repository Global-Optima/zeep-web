package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/categories"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/translations"
)

type CategoriesModule struct {
	*common.BaseModule
	Repo    categories.CategoryRepository
	Service categories.CategoryService
	Handler *categories.CategoryHandler
}

func NewCategoriesModule(base *common.BaseModule, auditService audit.AuditService, translationManager translations.TranslationManager) *CategoriesModule {
	repo := categories.NewCategoryRepository(base.DB)
	transactionManager := categories.NewTransactionManager(base.DB, repo, translationManager)
	service := categories.NewCategoryService(repo, transactionManager, base.Logger)
	handler := categories.NewCategoryHandler(service, auditService)

	base.Router.RegisterProductCategoriesRoutes(handler)

	return &CategoriesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
