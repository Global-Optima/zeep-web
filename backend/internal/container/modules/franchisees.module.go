package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
)

type FranchiseesModule struct {
	*common.BaseModule
	Repo    franchisees.FranchiseeRepository
	Service franchisees.FranchiseeService
	Handler *franchisees.FranchiseeHandler
}

func NewFranchiseesModule(base *common.BaseModule, auditService audit.AuditService) *FranchiseesModule {
	repo := franchisees.NewFranchiseeRepository(base.DB)
	service := franchisees.NewFranchiseeService(repo)
	handler := franchisees.NewFranchiseeHandler(service, auditService)

	base.Router.RegisterFranchiseeRoutes(handler)

	return &FranchiseesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
