package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/notifications"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/provisions/storeProvisions"
)

type ProvisionsModule struct {
	*common.BaseModule
	Repo                  provisions.ProvisionRepository
	Service               provisions.ProvisionService
	Handler               *provisions.ProvisionHandler
	StoreProvisionsModule *StoreProvisionsModule
}

func NewProvisionsModule(
	base *common.BaseModule,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	notificationService notifications.NotificationService,
) *ProvisionsModule {
	repo := provisions.NewProvisionRepository(base.DB)
	service := provisions.NewProvisionService(repo, base.Logger)
	handler := provisions.NewProvisionHandler(service, auditService, base.Logger)

	storeProvisionsModule := NewStoreProvisionsModule(base, service, franchiseeService, auditService, notificationService, repo)

	base.Router.RegisterProvisionsRoutes(handler)

	return &ProvisionsModule{
		BaseModule:            base,
		Repo:                  repo,
		Service:               service,
		Handler:               handler,
		StoreProvisionsModule: storeProvisionsModule,
	}
}

func NewStoreProvisionsModule(
	base *common.BaseModule,
	provisionService provisions.ProvisionService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	notificationService notifications.NotificationService,
	provisionRepo provisions.ProvisionRepository,
) *StoreProvisionsModule {
	repo := storeProvisions.NewStoreProvisionRepository(base.DB)
	service := storeProvisions.NewStoreProvisionService(repo, provisionRepo, notificationService, base.Logger)
	handler := storeProvisions.NewStoreProvisionHandler(service, provisionService, franchiseeService, auditService, base.Logger)

	base.Router.RegisterStoreProvisionsRoutes(handler)

	return &StoreProvisionsModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type StoreProvisionsModule struct {
	*common.BaseModule
	Repo    storeProvisions.StoreProvisionRepository
	Service storeProvisions.StoreProvisionService
	Handler *storeProvisions.StoreProvisionHandler
}
