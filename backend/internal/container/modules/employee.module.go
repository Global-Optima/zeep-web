package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
)

type EmployeesModule struct {
	*common.BaseModule
	Repo    employees.EmployeeRepository
	Service employees.EmployeeService
	Handler *employees.EmployeeHandler
}

func NewEmployeesModule(base *common.BaseModule) *EmployeesModule {
	repo := employees.NewEmployeeRepository(base.DB)
	service := employees.NewEmployeeService(repo, base.Logger)
	handler := employees.NewEmployeeHandler(service)

	base.Router.RegisterEmployeesRoutes(handler)
	base.Router.RegisterCommonEmployeesRoutes(handler)

	return &EmployeesModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
