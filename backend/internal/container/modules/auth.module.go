package modules

import (
	"github.com/Global-Optima/zeep-web/backend/internal/container/common"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/customers"
	"github.com/Global-Optima/zeep-web/backend/internal/modules/employees"
)

type AuthModule struct {
	*common.BaseModule
	Repo    auth.AuthenticationRepository
	Service auth.AuthenticationService
	Handler *auth.AuthenticationHandler
}

func NewAuthModule(
	base *common.BaseModule,
	customersRepo customers.CustomerRepository,
	employeesRepo employees.EmployeeRepository,
) *AuthModule {

	repo := auth.NewAuthenticationRepository(base.DB)
	service := auth.NewAuthenticationService(repo, customersRepo, employeesRepo, base.Logger)
	handler := auth.NewAuthenticationHandler(service)

	base.Router.RegisterAuthenticationRoutes(handler)

	return &AuthModule{
		BaseModule: base,
		Repo:       repo,
		Service:    service,
		Handler:    handler,
	}
}
