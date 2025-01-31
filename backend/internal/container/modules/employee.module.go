package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	storeEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees"
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

	storeEmployeeRepository := storeEmployees.NewStoreEmployeeRepository(base.DB)
	storeEmployeeService := storeEmployees.NewStoreEmployeeService(storeEmployeeRepository, repo, base.Logger)
	storeEmployeeHandler := storeEmployees.NewStoreEmployeeHandler(storeEmployeeService, service, franchiseeService, auditService)

	base.Router.RegisterEmployeesRoutes(handler)
	base.Router.RegisterEmployeeAccountRoutes(handler)
	base.Router.RegisterStoreEmployeeRoutes(storeEmployeeHandler)

	return &EmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
