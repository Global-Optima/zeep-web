package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/audit"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
	adminEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/adminEmployees"
	franchiseeEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/franchiseeEmployees"
	regionEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/regionEmployees"
	storeEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/storeEmployees"
	warehouseEmployees "github.com/Global-Optima/zeep-web/backend/internal/modules/employees/warehouseEmployees"
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

	storeEmployeesModule := NewStoreEmployeesModule(base, service, franchiseeService, auditService, repo)
	warehouseEmployeeModule := NewWarehouseEmployeesModule(base, service, regionService, auditService, repo)
	franchiseeEmployeesModule := NewFranchiseeEmployeesModule(base, service, franchiseeService, auditService, repo)
	regionEmployeesModule := NewRegionEmployeesModule(base, service, regionService, auditService, repo)
	adminEmployeesModule := NewAdminEmployeesModule(base, service, auditService, repo)

	base.Router.RegisterEmployeeAccountRoutes(
		storeEmployeesModule.Handler,
		warehouseEmployeeModule.Handler,
		franchiseeEmployeesModule.Handler,
		regionEmployeesModule.Handler,
		adminEmployeesModule.Handler,
	)

	base.Router.RegisterEmployeesRoutes(
		handler,
		storeEmployeesModule.Handler,
		warehouseEmployeeModule.Handler,
		franchiseeEmployeesModule.Handler,
		regionEmployeesModule.Handler,
		adminEmployeesModule.Handler,
	)

	return &EmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type StoreEmployeesModule struct {
	*common.BaseModule
	Repo    storeEmployees.StoreEmployeeRepository
	Service storeEmployees.StoreEmployeeService
	Handler *storeEmployees.StoreEmployeeHandler
}

func NewStoreEmployeesModule(
	base *common.BaseModule,
	employeeService employees.EmployeeService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	employeeRepo employees.EmployeeRepository,
) *StoreEmployeesModule {
	repo := storeEmployees.NewStoreEmployeeRepository(base.DB, employeeRepo)
	service := storeEmployees.NewStoreEmployeeService(repo, employeeRepo, base.Logger)
	handler := storeEmployees.NewStoreEmployeeHandler(service, employeeService, franchiseeService, auditService)

	return &StoreEmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type WarehouseEmployeesModule struct {
	*common.BaseModule
	Repo    warehouseEmployees.WarehouseEmployeeRepository
	Service warehouseEmployees.WarehouseEmployeeService
	Handler *warehouseEmployees.WarehouseEmployeeHandler
}

func NewWarehouseEmployeesModule(
	base *common.BaseModule,
	employeeService employees.EmployeeService,
	regionService regions.RegionService,
	auditService audit.AuditService,
	employeeRepo employees.EmployeeRepository,
) *WarehouseEmployeesModule {
	repo := warehouseEmployees.NewWarehouseEmployeeRepository(base.DB, employeeRepo)
	service := warehouseEmployees.NewWarehouseEmployeeService(repo, employeeRepo, base.Logger)
	handler := warehouseEmployees.NewWarehouseEmployeeHandler(service, employeeService, regionService, auditService)

	return &WarehouseEmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type FranchiseeEmployeesModule struct {
	*common.BaseModule
	Repo    franchiseeEmployees.FranchiseeEmployeeRepository
	Service franchiseeEmployees.FranchiseeEmployeeService
	Handler *franchiseeEmployees.FranchiseeEmployeeHandler
}

func NewFranchiseeEmployeesModule(
	base *common.BaseModule,
	employeeService employees.EmployeeService,
	franchiseeService franchisees.FranchiseeService,
	auditService audit.AuditService,
	employeeRepo employees.EmployeeRepository,
) *FranchiseeEmployeesModule {
	repo := franchiseeEmployees.NewFranchiseeEmployeeRepository(base.DB, employeeRepo)
	service := franchiseeEmployees.NewFranchiseeEmployeeService(repo, employeeRepo, base.Logger)
	handler := franchiseeEmployees.NewFranchiseeEmployeeHandler(service, employeeService, franchiseeService, auditService)

	return &FranchiseeEmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type RegionEmployeesModule struct {
	*common.BaseModule
	Repo    regionEmployees.RegionEmployeeRepository
	Service regionEmployees.RegionEmployeeService
	Handler *regionEmployees.RegionEmployeeHandler
}

func NewRegionEmployeesModule(
	base *common.BaseModule,
	employeeService employees.EmployeeService,
	regionService regions.RegionService,
	auditService audit.AuditService,
	employeeRepo employees.EmployeeRepository,
) *RegionEmployeesModule {
	repo := regionEmployees.NewRegionEmployeeRepository(base.DB, employeeRepo)
	service := regionEmployees.NewRegionEmployeeService(repo, employeeRepo, base.Logger)
	handler := regionEmployees.NewRegionEmployeeHandler(service, employeeService, regionService, auditService)

	return &RegionEmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}

type AdminEmployeesModule struct {
	*common.BaseModule
	Repo    adminEmployees.AdminEmployeeRepository
	Service adminEmployees.AdminEmployeeService
	Handler *adminEmployees.AdminEmployeeHandler
}

func NewAdminEmployeesModule(
	base *common.BaseModule,
	employeeService employees.EmployeeService,
	auditService audit.AuditService,
	employeeRepo employees.EmployeeRepository,
) *AdminEmployeesModule {
	repo := adminEmployees.NewAdminEmployeeRepository(base.DB)
	service := adminEmployees.NewAdminEmployeeService(repo, employeeRepo, base.Logger)
	handler := adminEmployees.NewAdminEmployeeHandler(service, employeeService, auditService)

	return &AdminEmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
