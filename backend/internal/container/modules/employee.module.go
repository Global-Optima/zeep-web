package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/franchisees"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/regions"
)

type EmployeesModule struct {
	*common.BaseModule
	Repo    employees.EmployeeRepository
	Service employees.EmployeeService
	Handler *employees.EmployeeHandler
}

func NewEmployeesModule(
	base *common.BaseModule,
	auditService audit.AuditService,
	franchiseeService franchisees.FranchiseeService,
	regionService regions.RegionService,
) *EmployeesModule {
	repo := employees.NewEmployeeRepository(base.DB)
	service := employees.NewEmployeeService(repo, base.Logger)
	handler := employees.NewEmployeeHandler(service, auditService, franchiseeService, regionService)

	base.Router.RegisterEmployeesRoutes(handler)
	base.Router.RegisterEmployeeAccountRoutes(handler)

	return &EmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
