package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/stores"
)

type StoresModule struct {
	*common.BaseModule
	Repo    stores.StoreRepository
	Service stores.StoreService
	Handler *stores.StoreHandler
}

func NewStoresModule(base *common.BaseModule, franchiseeService franchisees.FranchiseeService, auditService audit.AuditService) *StoresModule {
	repo := stores.NewStoreRepository(base.DB)
	service := stores.NewStoreService(repo, base.Logger)
	handler := stores.NewStoreHandler(service, franchiseeService, auditService)

	base.Router.RegisterStoresRoutes(handler)
	base.Router.RegisterCommonStoresRoutes(handler)

	return &StoresModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
